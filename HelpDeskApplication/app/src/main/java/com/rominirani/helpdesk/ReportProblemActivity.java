package com.rominirani.helpdesk;

import android.app.ProgressDialog;
import android.content.Intent;
import android.os.Handler;
import android.os.Message;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Spinner;
import android.widget.TextView;
import android.widget.Toast;

import org.json.JSONArray;
import org.json.JSONObject;

public class ReportProblemActivity extends AppCompatActivity {

    private String beaconname = "";
    private Button btnGetBeaconDetails;
    private Button btnReportIssue;
    private TextView tvBeaconName;
    private TextView txtAssetId;
    private TextView txtAssetName;
    private Spinner spinnerIssueType;
    private EditText edtIssueDetails;
    private EditText edtReportedBy;

    private ProgressDialog mProgressDialog = null;

    private Handler mHandler = new Handler() {
        public void handleMessage(Message msg) {
            mProgressDialog.dismiss();
            String action = msg.getData().getString("action");
            if (action.equalsIgnoreCase("GetBeaconDetails")) {
                String assetId = msg.getData().getString("assetId");
                String assetName = msg.getData().getString("assetName");
                //Toast.makeText(getBaseContext(), "Got Asset Details" + assetId + "," + assetName, Toast.LENGTH_SHORT).show();
                txtAssetId.setText(assetId);
                txtAssetName.setText(assetName);
            }
            else if (action.equalsIgnoreCase("ReportIssue")) {
                String ErrorCode = msg.getData().getString("ErrorCode");
                if (ErrorCode.equals("200")) {
                    //Clear the fields

                }
                else {

                }
            }
        }
    };

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_report_problem);

        txtAssetId = (TextView) findViewById(R.id.txtAssetId);
        txtAssetName = (TextView) findViewById(R.id.txtAssetName);
        spinnerIssueType = (Spinner) findViewById(R.id.spinnerIssueType);
        edtIssueDetails = (EditText) findViewById(R.id.edtIssueDetails);
        edtReportedBy = (EditText) findViewById(R.id.edtReportedBy);

        //Retrieve the extra parameters passed
        Intent i = getIntent();

        beaconname = i.getExtras().getString("beaconname");
        tvBeaconName = (TextView)findViewById(R.id.tvBeaconNameLabel);
        tvBeaconName.setText("Beacon Name: " + beaconname);

        //Set the On Click Listener for getting the Asset Details
        btnGetBeaconDetails = (Button)findViewById(R.id.btnGetBeaconDetails);
        btnGetBeaconDetails.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                //Make Network call to get the details
                mProgressDialog = ProgressDialog.show(ReportProblemActivity.this, "HelpDesk Application", "Searching ...", true, true);
                Thread t = new Thread(new Runnable() {

                    @Override
                    public void run() {
                        Bundle b = new Bundle();
                        b.putString("action","GetBeaconDetails");
                        String beaconData = MyApplication.doBeaconSearch(beaconname);
                        if (beaconData != null) {
                            //Parse the JSON
                            //Toast.makeText(getBaseContext(), "Got Asset Details JSON : " + beaconData, Toast.LENGTH_SHORT).show();
                            try {
                                JSONArray results = new JSONArray(beaconData);
                                JSONObject beacon = results.getJSONObject(0);
                                b.putString("assetId",beacon.getString("AssetId"));
                                b.putString("assetName",beacon.getString("AssetName"));
                            }
                            catch (Exception ex) {
                                b.putString("assetId","NA");
                                b.putString("assetName","NA");
                            }
                        }
                        else {
                            b.putString("assetId","NA");
                            b.putString("assetName","NA");
                        }
                        Message m = new Message();
                        m.setData(b);
                        mHandler.sendMessage(m);
                    }
                });
                t.start();
            }
        });

        //Set the On Click Listener for Raising the Issue
        btnReportIssue = (Button)findViewById(R.id.btnReportIssue);
        btnReportIssue.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                String issueDetails = edtIssueDetails.getText().toString();
                String issueRaisedBy = edtReportedBy.getText().toString();
                String issueType = spinnerIssueType.getSelectedItem().toString();

                if (issueDetails.trim().length() == 0) {
                    Toast.makeText(getBaseContext(), "Please enter Issue Details", Toast.LENGTH_SHORT).show();
                    return;
                }
                if (issueRaisedBy.trim().length() == 0) {
                    Toast.makeText(getBaseContext(), "Please enter your name", Toast.LENGTH_SHORT).show();
                    return;
                }
                if (issueType.equals("Select Issue Type")) {
                    Toast.makeText(getBaseContext(), "Please specify the type of issue", Toast.LENGTH_SHORT).show();
                    return;
                }

                //Make Network call to report the issue
                mProgressDialog = ProgressDialog.show(ReportProblemActivity.this, "HelpDesk Application", "Reporting the Issue ...", true, true);
                Thread t = new Thread(new Runnable() {

                    @Override
                    public void run() {
                        Bundle b = new Bundle();
                        b.putString("action","ReportIssue");
                        String raiseIssueData = MyApplication.doRaiseIssue(txtAssetId.getText().toString(),txtAssetName.getText().toString(),spinnerIssueType.getSelectedItem().toString(),edtIssueDetails.getText().toString(),edtReportedBy.getText().toString());
                        if (raiseIssueData != null) {
                            try {

                                JSONObject result = new JSONObject(raiseIssueData);
                                String ErrorCode = result.getString("ErrorCode");
                                b.putString("ErrorCode",ErrorCode);
                            }
                            catch (Exception ex) {
                                b.putString("ErrorCode",ex.getMessage());
                            }
                        }
                        else {
                            b.putString("ErrorCode","Error while making call to network");
                        }
                        Message m = new Message();
                        m.setData(b);
                        mHandler.sendMessage(m);
                    }
                });
                t.start();
            }
        });

    }
}
