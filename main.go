package main

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "os"
    "os/exec"
    "strings"
    "unicode"

    tea "github.com/charmbracelet/bubbletea"
    "gopkg.in/yaml.v3"
)

type Config map[string]string

type model struct {
    itemsAtInit  []string
    items  []string
    ptr   int
    selected map[int]struct{}
    isInSearch bool
    tabName []string
    tabItems [][]string
    tabIndex int
    tabItemIndex int
    searchQuerry string
}

func initialModel(input []string,) model {
    tab := get_tab()
    tabName := make([]string,0)
    tabName = append(tabName,"All")
    for k,_ := range tab{
        tabName = append(tabName,k)
    }
    return model{
        items:  input,
        itemsAtInit:  input,
        tabName: tabName,
        tabItems: get_tab_items(tab,input),
        selected: make(map[int]struct{}),
    }
}

func get_tab_items(tab map[string]string, items []string)[][]string{
    tab_item := make([][]string,len(tab)+1)
    i := 0
    tab_item[i] = items
    for _,v := range tab{
        i++
        tab_item[i]=filterItems(items,v)
        fmt.Println(tab_item[i])
    }
    return tab_item
}

func filterItems(items []string, querry string)[]string{
    new_items := make([]string,0)
    for _,item := range items{
        if strings.Contains(item,querry){
            new_items = append(new_items,item)
        }
    }
    return new_items
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if m.isInSearch{ 
        switch msg := msg.(type) {
        case tea.KeyMsg:
            switch msg.String() {
            case "backspace":
                if len(m.searchQuerry) > 0{
                    m.searchQuerry = m.searchQuerry[:len(m.searchQuerry)-1]
                }
            case "ctrl+c":
                return m, tea.Quit
            case "esc":
                m.isInSearch = false
            default:
                if len(msg.Runes)==1 && unicode.IsLetter(msg.Runes[0]){
                    m.searchQuerry += tea.KeyMsg.String(msg)
                }
            }
        }
        m.items = filterItems(m.itemsAtInit,m.searchQuerry)
        return m,nil
    }
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "s","/":
            if m.isInSearch{
                m.isInSearch = false
            }else{
                m.isInSearch = true
            }
        case "tab":
            m.items = m.tabItems[m.tabIndex] 
            m.tabIndex = (m.tabIndex + 1)%len(m.tabItems)
        case "d":
            m.items = append(m.items[:m.ptr],m.items[m.ptr+1:]...)
        case "ctrl+c", "esc":
            return m, tea.Quit
        case "up", "k":
            if m.ptr > 0 { m.ptr-- }
        case "down", "j":
            if m.ptr < len(m.items)-1 { m.ptr++ }
        case "enter", " ":
            _, ok := m.selected[m.ptr]
            if ok {
                delete(m.selected, m.ptr)
            } else {
                m.selected[m.ptr] = struct{}{}
            }
        }
    }
    return m, nil
}

func (m model) get_items_to_display() []string{
    nb_of_item_to_display := 6

    start := m.ptr - nb_of_item_to_display/2
    if start < 0{ start = 0 }

    end := m.ptr + nb_of_item_to_display/2
    if end > len(m.items){ end = len(m.items) }
    return m.items[start:end]
}

func (m model) View() string {
    s := fmt.Sprintf("-------\n")
    s += fmt.Sprintf("| %s |\n", m.tabName[m.tabIndex])
    s += fmt.Sprintf("-------\n")
    s += fmt.Sprintf("-%d-\n", m.ptr)

    key := "---"
    short := "---"
    long := "---"

    itemsToDisplay := m.get_items_to_display()
    start := m.ptr - len(itemsToDisplay)/2
    if start < 0 {
        start = 0
    }

    for displayIndex, choice := range itemsToDisplay {
        actualIndex := start + displayIndex
        if actualIndex >= len(m.items) {
            actualIndex = len(m.items) - 1
        }

        splited := strings.Split(choice, "{")
        if len(splited) >= 2 {
            key = splited[0]
            short = strings.Join(splited[1:], "")
            all_arg := strings.Split(short, ",")
            long = strings.Join(all_arg, "\n")
        }

        cursor := " "
        if actualIndex == m.ptr {
            cursor = "➡️"
        }

        if _, ok := m.selected[actualIndex]; ok {
            cursor = "⬇"
            s += fmt.Sprintf("%s  [%s] %s\n", cursor, key, long)
        } else {
            s += fmt.Sprintf("%s  [%s] %s\n", cursor, key, short)
        }
    }

    if m.isInSearch {
        s += fmt.Sprintf("is in search\n")
        s += fmt.Sprintf("< %s >\n", m.searchQuerry)
    }
    s += "\nPress esc to quit.\n"
    return s
}

func (m model) Init() tea.Cmd{
    return nil
}

func main() {
    get_tab()
    var reader *bufio.Reader

    fileInfo, _ := os.Stdin.Stat()
    if (fileInfo.Mode() & os.ModeNamedPipe) != 0 {
        reader = bufio.NewReader(os.Stdin)
    } else {
        reader = getClipboard()
    } 
    items := get_items(reader)
    if len(items)==0{
        items = get_items(reader)
    }

    p := tea.NewProgram(initialModel(items))
    if _, err := p.Run(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}

func getClipboard() *bufio.Reader {
    cmd := exec.Command("wl-paste")
    output, err := cmd.Output()
    if err != nil {
        return nil
    }
    return bufio.NewReader(bytes.NewReader(output))
}

func get_items(reader *bufio.Reader)[]string{
    var items []string
    for {
        line, err := reader.ReadString('\n')
        if err == io.EOF { break }
        if err != nil {
            fmt.Println("Error reading input:", err)
            os.Exit(1)
        }
        items = append(items, line)
    }
    return items
}

func get_tab()map[string]string{
    data, err := os.ReadFile("tab.yml")
    if err != nil {
        fmt.Println("warn: no tab file")
    }
    var config Config
    err = yaml.Unmarshal(data, &config)
    if err != nil {
        fmt.Println("err: deserialization")
    }
    return config
}
