# WAStats

WhatsApp group conversation statistics

v1.09 (the version starts of with this as there were earlier non-published versions going back a couple of years)

How to use?
Go to WhatsApp on your phone and select a group that you want to get conversation stats for. Select the menu option within that group and choose export data. You will be prompted to include or exclude media (the exact message may vary depending on the version of WhatsApp) and you should choose to exclude media. You will have to transfer this exported data to your computer. The easiest way is to pick the email option that WhatsApp will prompt you at the end of the export step and mail it to yourself. Or via Bluetooth directly to your computer.

You will only need to the file ending in ".txt" that contains conversations and you can ignore any other files present (usually .vcf files if contacts had been shared). Please note the name of this file.

WAStats -db exported_file_name -os iOS

The last parameter (iOS) will be Android if the data is from an Android phone.

WAStats runs and within a few seconds you will see a simple count of how many messages it has processed. It will also tell you that it has saved the results to a PDF file and it will give you the file name. Open the PDF and you will see the results - it shows total count of messages and senders, monthly and hourly distribution of messages and a list of senders and the number of messages they have sent in descending order.

It treats multiple short messages sent one after the other as just one message - this is based on how WhatsApp itself shows the messages on the screen.

There are a few more command line options, which will be shown when you run WAStats.exe without any options.
