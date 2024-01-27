# Up Bank FBAR Calculator

The US Department of the Treasury requires US citizens and residents to report their foreign financial accounts to the Treasury. This includes bank accounts, brokerage accounts, mutual funds, etc. The FBAR is a form that must be filed electronically with the Treasury. The FBAR is due on June 30th of each year, and covers the previous calendar year.

This tool is designed to help Up Bank customers calculate the maximum value of their Up Bank accounts for the FBAR. It is not an official tool, and is not endorsed by Up Bank, or (god forbid) the Department of the Treasury.

This tool is run via the command line, and must be passed an Up API token via the `UP_TOKEN` environment variable. To get an Up API token, follow the instructions [here](https://developer.up.com.au/#getting-started).

When run, the program will think get a list of transactions from the Up API, then collate them into per-account reports for every account you have with Up. It will then print out a short report for each account, and and create a CSV file for each account containing the transactions for that account. You should hold onto these CSVs for your record-keeping.

To run:
```Bash
UP_TOKEN=<your API token> YEAR=<the year you want to calculate the FBAR for> go run main.go
```

# Disclaimer!!!
I'm some third party rando. This software comes as is, with no warranty, etc, and i'm not liable for anything that happens to you or your money. I'm just some guy who made this software for to help with (sigh) filing my FBARs. This software is not endorsed by Up Bank, or the US Department of the Treasury, or anyone else, including me.

In using this software, you agree that you are using it at your own risk, and that you will not hold me liable for anything that happens to you or your money as a result of using this software. This software may produce incorrect data, and you should double check the data it produces before using it to file your FBAR.

This program and the information it produces are not financial advice. You should consult a professional before making any financial decisions, including (especially!) filing taxation information with the US Government.
