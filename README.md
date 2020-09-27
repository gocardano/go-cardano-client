
# go-cardano-client

_WORK IN PROGRESS_

This application is an attempt to integrate with the Cardano node via unix socket using golang.

The application uses it's own CBOR decoding/encoding implemention and the first feature implemented is querying the top block header hash.

## Usage

```
$ go-cardano-client -socket "/Users/macuser/Library/Application Support/Daedalus Mainnet/cardano-node.socket"
SlotNumber  :  2271537
Hash        :  a8bb6591131ee02984e94bd230b7d6eefc69101f9e8b16e2738fee67eda40857
BlockNumber :  2270047
```

```
$ go-cardano-client -socket "/Users/macuser/Library/Application Support/Daedalus Mainnet/cardano-node.socket" -debug
time="2020-10-23T17:08:15-05:00" level=trace msg="Multiplexed Request"
==========================================================================================
Header: Transmission Time: [1014223168], Mode: [0], Protocol ID: [0], Payload Length: [0]
------------------------------------------------------------------------------------------
Array: [2]
  PositiveInteger8(0)
  Map - Items: [3]
    - key: PositiveInteger32(1) / value: PositiveInteger32(764824073)
    - key: PositiveInteger32(2) / value: PositiveInteger32(764824073)
    - key: PositiveInteger32(3) / value: PositiveInteger32(764824073)
==========================================================================================
time="2020-10-23T17:08:15-05:00" level=debug msg="Attempting to write 29 bytes to socket"
time="2020-10-23T17:08:15-05:00" level=debug msg="Successfully written [29] bytes to socket"
time="2020-10-23T17:08:15-05:00" level=debug msg="Received response header" header="0x57 0x8f 0x40 0x87 0x80 0x00 0x00 0x08 " payloadLength=8 protocolID=0 transmissionTime=1469005959
time="2020-10-23T17:08:15-05:00" level=debug msg="Read packet" readCount=8 totalReadCount=8
time="2020-10-23T17:08:15-05:00" level=debug msg="Successfully read 8 bytes from socket" response="0x83 0x01 0x01 0x1a 0x2d 0x96 0x4a 0x09 "
time="2020-10-23T17:08:15-05:00" level=trace msg="Starting to iterate on array with length: 3"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array PositiveInteger8(1)"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array PositiveInteger8(1)"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array PositiveInteger32(764824073)"
time="2020-10-23T17:08:15-05:00" level=trace msg="Array of [3] length reached [3] items"
time="2020-10-23T17:08:15-05:00" level=info msg="Shelley container has [1] CBOR encoded data items"
time="2020-10-23T17:08:15-05:00" level=trace msg="Multiplexed Response"
==========================================================================================
Header: Transmission Time: [1469005959], Mode: [1], Protocol ID: [0], Payload Length: [8]
------------------------------------------------------------------------------------------
Array: [3]
  PositiveInteger8(1)
  PositiveInteger8(1)
  PositiveInteger32(764824073)
==========================================================================================
time="2020-10-23T17:08:15-05:00" level=debug msg="Handshake was successful" extraParams=764824073 versionNumber=1
time="2020-10-23T17:08:15-05:00" level=trace msg="Multiplexed Request"
==========================================================================================
Header: Transmission Time: [1014667168], Mode: [0], Protocol ID: [5], Payload Length: [0]
------------------------------------------------------------------------------------------
Array: [1]
  PositiveInteger8(0)
==========================================================================================
time="2020-10-23T17:08:15-05:00" level=debug msg="Attempting to write 10 bytes to socket"
time="2020-10-23T17:08:15-05:00" level=debug msg="Successfully written [10] bytes to socket"
time="2020-10-23T17:08:15-05:00" level=debug msg="Received response header" header="0x57 0x8f 0x58 0x5a 0x80 0x05 0x00 0x31 " payloadLength=49 protocolID=5 transmissionTime=1469012058
time="2020-10-23T17:08:15-05:00" level=debug msg="Read packet" readCount=49 totalReadCount=49
time="2020-10-23T17:08:15-05:00" level=debug msg="Successfully read 49 bytes from socket" response="0x83 0x03 0x80 0x82 0x82 0x1a 0x00 0x24 0x00 0xa9 0x58 0x20 0x73 0x46 0x07 0x60 0x8c 0x70 0xe0 0x70 0x57 0x8f 0x63 0x3f 0xf2 0xf5 0xfa 0xa8 0xee 0x58 0x83 0x1b 0xdb 0x0d 0xff 0xb2 0x97 0x9b 0xfb 0x47 0x65 0x9f 0x1f 0x8a 0x1a 0x00 0x23 0xfa 0xcf "
time="2020-10-23T17:08:15-05:00" level=trace msg="Starting to iterate on array with length: 3"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array PositiveInteger8(3)"
time="2020-10-23T17:08:15-05:00" level=trace msg="Starting to iterate on array with length: 0"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array Array: [0]"
time="2020-10-23T17:08:15-05:00" level=trace msg="Starting to iterate on array with length: 2"
time="2020-10-23T17:08:15-05:00" level=trace msg="Starting to iterate on array with length: 2"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array PositiveInteger32(2359465)"
time="2020-10-23T17:08:15-05:00" level=trace msg="Reading bytes of payload length: 32"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array ByteString - Length: [32]; Value: [734607608c70e070578f633ff2f5faa8ee58831bdb0dffb2979bfb47659f1f8a];"
time="2020-10-23T17:08:15-05:00" level=trace msg="Array of [2] length reached [2] items"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array Array: [2]"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array PositiveInteger32(2357967)"
time="2020-10-23T17:08:15-05:00" level=trace msg="Array of [2] length reached [2] items"
time="2020-10-23T17:08:15-05:00" level=trace msg="Found another item in the array Array: [2]"
time="2020-10-23T17:08:15-05:00" level=trace msg="Array of [3] length reached [3] items"
time="2020-10-23T17:08:15-05:00" level=info msg="Shelley container has [1] CBOR encoded data items"
time="2020-10-23T17:08:15-05:00" level=trace msg="Multiplexed Response"
==========================================================================================
Header: Transmission Time: [1469012058], Mode: [1], Protocol ID: [5], Payload Length: [49]
------------------------------------------------------------------------------------------
Array: [3]
  PositiveInteger8(3)
  Array: [0]
  Array: [2]
    Array: [2]
      PositiveInteger32(2359465)
      ByteString - Length: [32]; Value: [734607608c70e070578f633ff2f5faa8ee58831bdb0dffb2979bfb47659f1f8a];
    PositiveInteger32(2357967)
==========================================================================================
SlotNumber  :  2359465
Hash        :  734607608c70e070578f633ff2f5faa8ee58831bdb0dffb2979bfb47659f1f8a
BlockNumber :  2357967
```


# Developer Notes

## How to Capture Packets

Since it is not anymore possible to capture traffic via the Unix socket, the workaround would be to proxy a fake unix socket to a TCP port (where the tcpdump will be captured).

```
fake-unix-socket -> TCP PORT 6000 -> cardano-unix-socket
                    (tcpdump here)  

# Step 1: Create TCP port to cardano-unix-socket
socat TCP-LISTEN:6000,reuseaddr,fork UNIX-CONNECT:node.socket

# Step 2: Create fake unix socket where client would connect to*
socat UNIX-LISTEN:fake.socket,fork TCP-CONNECT:127.0.0.1:6000

# Step 3: Packet capture on the TCP port
tcpdump -i lo -f 'tcp port 6000' -w dump.pcap
```

* cardano-node connects to unix socket as specified by the env var `CARDANO_NODE_SOCKET`
```
export CARDANO_NODE_SOCKET_PATH=/path/to/unix.socket
```

Reference: https://unix.stackexchange.com/questions/219853/how-to-passively-capture-from-unix-domain-sockets-af-unix-socket-monitoring


## Analyzing Packets

### Handshake Request

```
# Raw Packets:
0000   00 00 00 00 00 00 00 00 00 00 00 00 08 00 45 00   ..............E.
0010   00 4d d1 c1 40 00 40 06 6a e7 7f 00 00 01 7f 00   .M..@.@.j.......
0020   00 01 e9 d2 17 70 b2 bc 85 81 d3 a4 a8 78 80 18   .....p.......x..
0030   01 56 fe 41 00 00 01 01 08 0a 74 19 ad 4b 74 19   .V.A......t..Kt.
0040   ad 3e c1 70 97 d5 00 00 00 11 82 00 a2 01 1a 2d   .>.p...........-
0050   96 4a 09 19 80 02 1a 2d 96 4a 09                  .J.....-.J.

----

Transmission time : c1 70 97 d5 
Mux Initiator     : 00 00
Payload length    : 00 11 (Decimal: 17)
Payload           : 82 00 a2 01 1a 2d 96 4a 09 19 80 02 1a 2d 96 4a 09 

----

Decoded CBOR Payload (via cbor.me):
Value: [0, {1: 764824073, 32770: 764824073}]

82                # array(2)
   00             # unsigned(0)
   A2             # map(2)
      01          # unsigned(1)
      1A 2D964A09 # unsigned(764824073)
      19 8002     # unsigned(32770)
      1A 2D964A09 # unsigned(764824073)
```

### Handshake Response

```
# Raw Packets:
0000   00 00 00 00 00 00 00 00 00 00 00 00 08 00 45 00   ..............E.
0010   00 46 0b 7c 40 00 40 06 31 34 7f 00 00 01 7f 00   .F.|@.@.14......
0020   00 01 17 70 e9 d2 d3 a4 a8 78 b2 bc 85 9a 80 18   ...p.....x......
0030   01 56 fe 3a 00 00 01 01 08 0a 74 19 ad 4b 74 19   .V.:......t..Kt.
0040   ad 4b c1 70 b3 81 80 00 00 0a 83 01 19 80 02 1a   .K.p............
0050   2d 96 4a 09                                       -.J.

----

Transmission time : c1 70 b3 81
Mux Initiator     : 80 00
Payload length    : 00 0a (Decimal: 10)
Payload           : 83 01 19 80 02 1a 2d 96 4a 09

----

Decoded CBOR Payload (via cbor.me):
Value: [1, 32770, 764824073]

83             # array(3)
   01          # unsigned(1)
   19 8002     # unsigned(32770)
   1A 2D964A09 # unsigned(764824073)
```


# References

* [Ouroboros Wireshark Plugin](https://github.com/input-output-hk/ouroboros-network/tree/02ff4eaeb8c4fa47a6b6d99cbd1219f209de3c05/ouroboros-network/wireshark-plugin)
* [Shelley Network Specs](https://hydra.iohk.io/build/4110312/download/2/network-spec.pdf)

