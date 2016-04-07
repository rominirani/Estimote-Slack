# Estimote-Slack

This is a Help Desk application that uses Estimote Beacons to make the task of reporting issues easier. The issues are directly reported into the Slack Channel, where the team can get immediately notified.

Check out the [**EstimoteSlackIntegration.pdf**](https://github.com/rominirani/Estimote-Slack/blob/master/EstimoteSlackIntegration.pdf) document for Application Flow and screenshots

There are two applications in this project:

1. Android Application : This application (**HelpDeskApplication**) detects the estimote beacon. Each Estimote Beacon is configured with a Hardware Asset, for e.g. a Printer. So, if there is an issue with the printer, the user pulls out the application which detects that the Estimote Beacon is present and it pops up an Issue Reporting screen in the Android application. The issue is then filled out and sent to the Server.
2. The Server application (**helpdesk**) is a Google App Engine application that receives the issue reported by the Android application. The issue is saved locally in the datastore and a Slack Notification is raised in the #general Slack Channel for that team. 

# Hardware
You will need Estimote Beacons, atleast one for this application

#Setup
The Server application exposes an API that is used to configure the beacons first. These beacons have to be configured first and associated with the hardware devices. The Slack Channel Incoming Notifications Webhook channel needs to be configured too. In the Android application, you will need to setup the Beacon Ids that you are going to use for the application.

Let us do this step by step:
## Server side Configuration

1. Clear All Data
[http://helpdeskslack.appspot.com/clearData](http://helpdeskslack.appspot.com/clearData)

2. Configure All Assets with Beacons. 

Assuming that you have a beacon or two in your Estimote Cloud account, you will need to invoke the /addAsset URL endpoint for each Asset that you plan to add.



    helpdeskslack.appspot.com/addAsset

Each Asset to be added requires the following request parameters.



- AssetId=`<assetid>`
- AssetName=`<assetname>`
- BeaconId=`<beaconid>`     #This is your Estimote Beacon Id set in the Estimote Cloud
- BeaconName=`<beaconname>` #This is your Estimote Beacon Name set in the Estimote Cloud
- AppType=Asset           #Currently set to Asset
- AppData=                #Currently set to Empty
- Status=ACTIVE

For example:

[> http://helpdeskslack.appspot.com/addAsset?AssetId=PRINTER-1&AssetName=HPLaserJetPlus&BeaconId=B9407F30...&BeaconName=B1&AppType=Asset&AppData=&Status=ACTIVE](http://helpdeskslack.appspot.com/addAsset?AssetId=PRINTER-1&AssetName=HPLaserJetPlus&BeaconId=B9407F30...&BeaconName=B1&AppType=Asset&AppData=&Status=ACTIVE)

The next step is to configure the Slack Webhook URL. You will need to note down your Slack Incoming Webhook URL and then invoke the `/addSlackTeamWebhook` endpoint to configure it in the Server application. 

`helpdeskslack.appspot.com/addSlackTeamWebhook?webhookurl=<Your_Slack_Incoming_Webhook_URL>
`

You can check the configured value of Slack Webhook URL:

 `http://helpdeskslack.appspot.com/SlackTeamWebhook`

## Android App Configuration
Import the Android application project : **HelpDeskApplication**.

Go to the following files: 


- `com.rominirani.helpdesk.MainActivity.java`, go to line 43 and replace `YOUR_BEACON_ID` with one of the 
Beacon Ids that you have configured in above step.
- Build the Android Project again and run it.
- Follow the steps as described in the [**EstimoteSlackIntegration.pdf**](https://github.com/rominirani/Estimote-Slack/blob/master/EstimoteSlackIntegration.pdf) document
