# envoy-pvoutput

Reads data from the Envoy S usage consumption meter using the (undocumented) production.json interface and posts this to PVOutput.org.

## Example

    envoy-pvoutput -ENVOYIP 192.168.1.1 -PVOUTPUTAPIKEY abcdefghijklmnopqrstuvwxyz -PVOUTPUTSYSTEMID 12345

## Usage:

    envoy-pvoutput
    
    -ENVOYIP string
        IP Address of Envoy S to retrieve data from
    -ENVOYPORT int
        Port of the Envoy S to retrieve data from (default 80)
    -POLLINTERVALSECONDS int
        Polling interval in seconds (default 300)
    -PVOUTPUTAPIKEY string
        PVOutput.org API Key to use to post data
    -PVOUTPUTSYSTEMID int
        PVOutput.org system ID for the Envoy S
    -TIMEZONE string
        Timezone of the Envoy S. If unset, same as current local timezone
