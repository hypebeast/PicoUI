package org.sebastianruml.miuiapp.helpers;


public class Utils {
	private static final int hours_in_day = 24;
	private static final int mins_in_hour = 60;
	private static final int secs_in_min = 60;
	
	static public String ConvertSecondsToTimeString(String secondsString) {
		StringBuilder timeString = new StringBuilder();
		
		int totalSeconds = 0;
		try {
			float s = Float.parseFloat(secondsString);
			totalSeconds = (int)Math.round(s);
		} catch (NumberFormatException e) {
			
		}
		
		int seconds = totalSeconds % secs_in_min;
		int minutes = totalSeconds / secs_in_min % mins_in_hour;
		int hours = totalSeconds / secs_in_min / mins_in_hour % hours_in_day;
		int days = totalSeconds / secs_in_min / mins_in_hour / hours_in_day;
		
		if (days > 0) {
			timeString.append(String.valueOf(days));
			timeString.append(" d ");
		}
		
		if (hours > 0) {
			timeString.append(String.valueOf(hours));
			timeString.append(" h ");
		}
		
		if (minutes > 0) {
			timeString.append(String.valueOf(minutes));
			timeString.append(" m ");
		}

		timeString.append(String.valueOf(seconds));
		timeString.append(" s");
		
		return timeString.toString();
	}
}
