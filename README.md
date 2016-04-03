# Estimote-Slack

This is a Help Desk application that uses Estimote Beacons to make the task of reporting issues easier. The issues are directly reported into the Slack Channel, where the team can get immediately notified.

Check out the **EstimoteSlackIntegration.pdf** document for Application Flow and screenshots

There are two applications in this project:

1. Android Application : This application (**HelpDeskApplication**) detects the estimote beacon. Each Estimote Beacon is configured with a Hardware Asset, for e.g. a Printer. So, if there is an issue with the printer, the user pulls out the application which detects that the Estimote Beacon is present and it pops up an Issue Reporting screen in the Android application. The issue is then filled out and sent to the Server.
2. The Server application (**helpdesk**) is a Google App Engine application that receives the issue reported by the Android application. The issue is saved locally in the datastore and a Slack Notification is raised in the #general Slack Channel for that team. 

# Hardware
You will need Estimote Beacons, atleast one for this application

#Setup
The Server application exposes an API that is used to configure the beacons first. These beacons have to be configured first and associated with the hardware devices. The Slack Channel Incoming Notifications Webhook channel needs to be configured too. In the Android application, you will need to setup the Beacon Ids that you are going to use for the application.
