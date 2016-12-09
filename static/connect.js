
var config ={"iceServers":[{"urls":["stun:stun.l.google.com:19302"]}]};

// Chrome / Firefox compatibility
window.PeerConnection = window.RTCPeerConnection || window.mozRTCPeerConnection || window.webkitRTCPeerConnection;
window.RTCIceCandidate = window.RTCIceCandidate || window.mozRTCIceCandidate;
window.RTCSessionDescription = window.RTCSessionDescription || window.mozRTCSessionDescription;


function fromWebSocket(addr, cb) {
    const ws = new WebSocket(addr);
    const pc = new PeerConnection(config, {
        optional: [
            { DtlsSrtpKeyAgreement: true },
            { RtpDataChannels: false },
        ],
    });

    ws.onopen = function(evt) { console.log("ws: open", evt); }
    ws.onclose = function(evt) { console.log("ws: close", evt); }

    ws.onmessage = function(msg) {
        console.log("ws: msg", msg)
        const sdp = new RTCSessionDescription(JSON.parse(msg.data));
        const err = pc.setRemoteDescription(sdp);
        console.log("ws: SDP " + sdp.type + " successfully received.");

        pc.createAnswer().then(function(answer) {
            pc.setLocalDescription(answer);
            ws.send(JSON.stringify(pc.localDescription));
        });
    };

    pc.onicecandidate = function(evt) {
        const candidate = evt.candidate;
        // Chrome sends a null candidate once the ICE gathering phase completes.
        // In this case, it makes sense to send one copy-paste blob.
        if (null == candidate) {
            return;
        }
    }
    pc.ondatachannel = function(dc) {
        console.log("Data Channel established: ", dc);
        channelReady(dc.channel);
    }

    function channelReady(channel) {
        // Creating the first data channel triggers ICE negotiation.
        channel.onopen = function() {
            console.log("Data channel opened!");
        }
        channel.onclose = function() {
            console.log("Data channel closed.");
        }
        channel.onerror = function() {
            console.log("Data channel error!!", this.arguments);
        }
        channel.onmessage = function(msg) {
            console.log("channel: ", msg);
        }

        cb(channel);
    };

    const channel = pc.createDataChannel("test");

    return pc;
}
