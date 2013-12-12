package org.sebastianruml.miuiapp.fragments;

import org.sebastianruml.miuiapp.PicoUiStatus;
import org.sebastianruml.miuiapp.R;
import org.sebastianruml.miuiapp.SystemControlTask;

import android.app.AlertDialog;
import android.app.Fragment;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.TextView;
import android.content.DialogInterface;

public class HomeFragment extends Fragment {
	private TextView statusText;
	private TextView uptimeText;
	private TextView loadText;
	private TextView cpuText;
	private TextView memText;
	private Button rebootButton;
	private Button shutdownButton;
	private String mHostUrl;
	
	public HomeFragment(){}
    
    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
            Bundle savedInstanceState) {
  
    	final View rootView = inflater.inflate(R.layout.fragment_home, container, false);
    	
    	Bundle bundle = getArguments();
    	mHostUrl = bundle.getString("hostUrl");
		
		statusText = (TextView)rootView.findViewById(R.id.textStatus);
		statusText.setText("Connecting...");
		uptimeText = (TextView)rootView.findViewById(R.id.textUptime);
		uptimeText.setText("-");
		loadText = (TextView)rootView.findViewById(R.id.textLoad);
		loadText.setText("-");
		cpuText = (TextView)rootView.findViewById(R.id.textCpu);
		cpuText.setText("-");
		memText = (TextView)rootView.findViewById(R.id.textMem);
		memText.setText("-");
		rebootButton = (Button)rootView.findViewById(R.id.buttonReboot);
		rebootButton.setOnClickListener(new View.OnClickListener() {
			
			@Override
			public void onClick(View arg0) {
				AlertDialog.Builder builder = new AlertDialog.Builder(rootView.getContext());
				builder.setMessage("Do you really want to reboot your system?")
				       .setCancelable(true)
				       .setPositiveButton("OK", new DialogInterface.OnClickListener() {
				           public void onClick(DialogInterface dialog, int id) {
				        	   SystemControlTask task = new SystemControlTask();
				        	   task.execute("reboot", mHostUrl);
				           }
				       })
				       .setNegativeButton("Cancel", null);
				AlertDialog alert = builder.create();
				alert.show();
			}
		});
		shutdownButton = (Button)rootView.findViewById(R.id.buttonShutdown);
		shutdownButton.setOnClickListener(new View.OnClickListener() {
			
			@Override
			public void onClick(View v) {
				AlertDialog.Builder builder = new AlertDialog.Builder(rootView.getContext());
				builder.setMessage("Do you really want to shutdown your system?")
				       .setCancelable(true)
				       .setPositiveButton("OK", new DialogInterface.OnClickListener() {
				           public void onClick(DialogInterface dialog, int id) {
				        	   SystemControlTask task = new SystemControlTask();
				        	   task.execute("shutdown", mHostUrl);
				           }
				       })
				       .setNegativeButton("Cancel", null);
				AlertDialog alert = builder.create();
				alert.show();
			}
		});
          
        return rootView;
    }
    
    public void SetStatusText(String text) {
    	statusText.setText(text);
    }
    
    public void SetStatus(PicoUiStatus status) {
		uptimeText.setText(status.GetUptime());
		loadText.setText(status.GetLoad());
		memText.setText(Integer.toString(status.GetFreeMem()) + " MB");
		cpuText.setText(Integer.toString(status.GetCpuLoad()) +  " %");
    }
}
