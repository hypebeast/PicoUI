package org.sebastianruml.miuiapp;

public class PicoUiStatus {
	private String connectionStatus;
	private String uptime;
	private String load;
	private int freeMem;
	private int cpuLoad;
	
	public String GetConnectionStatus() {
		return connectionStatus;
	}
	
	public void SetConnectionStatus(String status) {
		connectionStatus = status;
	}
	
	public String GetUptime() {
		return uptime;
	}
	
	public void SetUptime(String uptime) {
		this.uptime = uptime;
	}
	
	public String GetLoad() {
		return load;
	}
	
	public void SetLoad(String load) {
		this.load = load;
	}
	
	public int GetFreeMem() {
		return freeMem;
	}
	
	public void SetFreeMem(int mem) {
		this.freeMem = mem;
	}
	
	public int GetCpuLoad() {
		return this.cpuLoad;
	}
	
	public void SetCpuLoad(int cpu) {
		this.cpuLoad = cpu;
	}
}
