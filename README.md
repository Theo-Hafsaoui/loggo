# Loggo

<p align="center">
  <img src="https://github.com/user-attachments/assets/71c2171d-b15f-46ca-a451-23fd0d262685" alt="loggo" width="600">
</p>

Loggo is a simple log parser designed for quick exploration of log dumps. It allows navigation, filtering, and searching through logs efficiently.

It's *very* specifique to my use case and *not* very pretty or well designed, and it has many flaws.

## Installation

Ensure you have Go installed, then build the project:

```sh
make build
```

## Usage

You can pipe logs into Loggo for parsing:

```sh
cat ./example.log | ./loggo
```

Or if no parameter use the clipboard.  
Loggo supports interactive navigation, searching, and tab-based filtering to help you analyze logs quickly.

