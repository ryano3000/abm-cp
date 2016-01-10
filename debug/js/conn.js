var wsocket = new WebSocket('ws://localhost:8080/ws')

wsocket.onopen = function () {
  console.log('websocket connection init')
  var viz = new p5(sketch, 'abm-viewport')
}

wsocket.onmessage = function (e) {
  var rawmsg = JSON.parse(e.data)
  console.log(rawmsg)
  if (rawmsg.type === 'render') {
    drawlist.cpp = rawmsg.data.cpp
    drawlist.vp = rawmsg.data.vp
    drawlist.bg = rawmsg.data.bg
  }
}

wsocket.onerror = function (e) {
  console.log('ERROR: ' + e)
}
