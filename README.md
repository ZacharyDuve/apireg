# go-distributed-api-registry
Distributed Api Registry written in go. Currently fully functional though still in BETA

# Function:
A simple leaderless in memory only Api Registry. The idea is that every instance of the registry keeps a complete list of all Apis that it knows about. Each registry publishes and listens over multicast for packets containing API registry information. While not the best networks with a large number of deployments I needed something on my local home network which would allow for simple auto discovery of APIs along with enough information to be able to connect. The Registry does its best to keep records up to date but it isn't guarrentied that a record is still active so it up to the code that actually connects to handle nothing listening anymore. 

# Current Configs:
Current config which is subject to change is packets are sent for update at least every 30 seconds and retired after 90 seconds if no packet for update has been received

Current Multicast config is IP of "224.0.0.78" and port of 5324

# What an API is:
An API is simply a Name, Version, and Port that you have your API setup for.
    All registration packets are encoded into JSON there currently is a soft limit of a packet containing 1200 bytes

# Functions available:
Registry has 3 functions:

    RegisterApi(name string, version string, port int) error

Which is used for registering your current applications API. You can register as many unique sets of APIs within a registry as you want.

    GetAvailableApis() []Api

Which returns every API that the registry knows about and is still tracking

    GetApisByApiName(name string) []Api

Which returns all APIs that the registry knows about and is tracking for a given name only. Will return multiple entries if version, ip, or port differs

Note: There are no functions currently to remove or delete a registration in a registry. I didn't think that they were needed as I currently only see adding on bootup and then using the lookup feature. If there would be changes to my published APIs then whole app would be brought down first which would completely reset the registry

# Example usage:
For my current model railroad I have multiple switch machine driver servers. Each would say publish "Name: SMDS, Version: v1, Port: 80". I also would have a single 'Turnout Central Command' server who would be able to talk to SMDS servers of v1. The registry allows for the 'Turnout Central Command' server to identify which IPs have SMDS v1 running along with the port. Then from there SMDS client software can connect to each server without having to know hostnames or IPs from a manual config.
