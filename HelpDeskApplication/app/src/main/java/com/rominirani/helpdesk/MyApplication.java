package com.rominirani.helpdesk;

import android.app.Application;

import com.estimote.sdk.EstimoteSDK;

import java.io.BufferedReader;
import java.io.InputStreamReader;
import java.net.URLEncoder;

import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.Response;

public class MyApplication extends Application {

    @Override
    public void onCreate() {
        super.onCreate();

        EstimoteSDK.initialize(getApplicationContext(), "help-desk-application-f2w", "APP_ID");

        // uncomment to enable debug-level logging
        // it's usually only a good idea when troubleshooting issues with the Estimote SDK
//        EstimoteSDK.enableDebugLogging(true);
    }

    public static String doBeaconSearch(String beaconName) {
        String url = "http://YOUR_APP/assets?BeaconName="
                + URLEncoder.encode(beaconName);

        OkHttpClient client = new OkHttpClient();
        Request request = new Request.Builder()
                .url(url)
                .build();
        try {
            Response response = client.newCall(request).execute();
            return response.body().string();
        }
        catch (Exception ex) {
            return null;
        }

    }

    public static String doRaiseIssue(String assetId, String assetName, String issueType, String issueDetails, String issueRaisedBy) {
        String url = "http://YOUR_APP/raiseIssue?"
                + "AssetId=" + URLEncoder.encode(assetId)
                + "&AssetName=" + URLEncoder.encode(assetName)
                + "&IssueType=" + URLEncoder.encode(issueType)
                + "&IssueDetails=" + URLEncoder.encode(issueDetails)
                + "&IssueRaisedBy=" + URLEncoder.encode(issueRaisedBy)
                + "&Status=ACTIVE";

        OkHttpClient client = new OkHttpClient();
        Request request = new Request.Builder()
                .url(url)
                .build();
        try {
            Response response = client.newCall(request).execute();
            return response.body().string();
        }
        catch (Exception ex) {
            return null;
        }

    }
}
