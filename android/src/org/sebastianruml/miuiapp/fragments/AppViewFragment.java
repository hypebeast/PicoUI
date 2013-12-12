package org.sebastianruml.miuiapp.fragments;

import org.sebastianruml.miuiapp.R;

import android.annotation.SuppressLint;
import android.app.Fragment;
import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.webkit.WebView;

public class AppViewFragment extends Fragment {
	private WebView mWebView;
	private String mAppUrl;
	
	public AppViewFragment () {}
	
	@SuppressLint("SetJavaScriptEnabled")
	@Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
            Bundle savedInstanceState) {
  
    	View rootView = inflater.inflate(R.layout.fragment_appview, container, false);
    	
    	mWebView = (WebView) rootView.findViewById(R.id.webView);
    	mWebView.getSettings().setJavaScriptEnabled(true);
    	
    	//mWebView.loadUrl("file:///android_assets/loading.html");
    	
    	Bundle bundle = getArguments();
    	mAppUrl = bundle.getString("hostUrl");
    	
    	mWebView.loadUrl(mAppUrl);
    	
        return rootView;
    }
}
