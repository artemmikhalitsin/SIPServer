# SIPServer
Server which holds SIP registrations and allows a client to query registrations

## Usage

    cd app
    go build
    ./SIPServer

By default it listens on port 4444 and looks for registration file "regs" in the local folder. This can be modified with environment variables, i.e.

    SIP_PORT=4000 SIP_REGS_FILE=/usr/local/registrations ./SIPServer

Alternatively, you can use the provided docker-compose configuration:

    docker-compose up

Then connect on the specified port using TCP connection utility

    nc localhost 4444

After this, you can look up records by their addressOfRecord

### Example
      $ ~/SIPServer$ nc localhost 4444
      014e9cc9ea34446a2b000100620005
      {"AddressOfRecord":"014e9cc9ea34446a2b000100620005","TenantId":"0127d974-f9f3-0704-2dee-000100420001","Uri":"sip:014e9cc9ea34446a2b000100620005@193.68.97.166;jbcuser=cpe70","Contact":"\u003csip:014e9cc9ea34446a2b000100620005@65.237.236.107;jbcuser=cpe70\u003e;methods=\"INVITE, ACK, BYE, CANCEL, OPTIONS, INFO, MESSAGE, SUBSCRIBE, NOTIFY, PRACK, UPDATE, REFER\"","Path":["\u003csip:Mi0xOTkuMTkyLjE2NS4xOTQtMTk2MjI@74.144.114.94:5060;lr\u003e"],"Source":"45.244.97.230:19622","Target":"192.250.34.27:5061","UserAgent":"polycom.soundstationip.5000","RawUserAgent":"PolycomSoundStationIP-SSIP_5000-UA/199.152.0.142","Created":"2016-12-12T22:42:04.538Z","LineId":"0148b27e-1149-a336-b632-000100620005"}
      014e9cc9ea34446a2b000100620005
      {"AddressOfRecord":"014e9cc9ea34446a2b000100620005","TenantId":"0127d974-f9f3-0704-2dee-000100420001","Uri":"sip:014e9cc9ea34446a2b000100620005@193.68.97.166;jbcuser=cpe70","Contact":"\u003csip:014e9cc9ea34446a2b000100620005@65.237.236.107;jbcuser=cpe70\u003e;methods=\"INVITE, ACK, BYE, CANCEL, OPTIONS, INFO, MESSAGE, SUBSCRIBE, NOTIFY, PRACK, UPDATE, REFER\"","Path":["\u003csip:Mi0xOTkuMTkyLjE2NS4xOTQtMTk2MjI@74.144.114.94:5060;lr\u003e"],"Source":"45.244.97.230:19622","Target":"192.250.34.27:5061","UserAgent":"polycom.soundstationip.5000","RawUserAgent":"PolycomSoundStationIP-SSIP_5000-UA/199.152.0.142","Created":"2016-12-12T22:42:04.538Z","LineId":"0148b27e-1149-a336-b632-000100620005"}
      82381
      Could not find registration with address 82381
      Connection closed.
