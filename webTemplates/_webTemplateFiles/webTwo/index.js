'use strict' //Enables strict mode is JavaScript
let url = new URL(document.location.href)
let gc_region = url.searchParams.get('gc_region')
let gc_clientId = url.searchParams.get('gc_clientId')
let gc_redirectUrl = url.searchParams.get('gc_redirectUrl')
let userId

//Getting and setting the GC details from dynamic URL and session storage
gc_region ? sessionStorage.setItem('gc_region', gc_region) : (gc_region = sessionStorage.getItem('gc_region'))
gc_clientId ? sessionStorage.setItem('gc_clientId', gc_clientId) : (gc_clientId = sessionStorage.getItem('gc_clientId'))
gc_redirectUrl ? sessionStorage.setItem('gc_redirectUrl', gc_redirectUrl) : (gc_redirectUrl = sessionStorage.getItem('gc_redirectUrl'))

let platformClient = require('platformClient')
const client = platformClient.ApiClient.instance
const uapi = new platformClient.UsersApi()
const napi = new platformClient.NotificationsApi()

async function start() {
  try {
    client.setEnvironment(gc_region)
    client.setPersistSettings(true, '_mm_')

    console.log('%cLogging in to Genesys Cloud', 'color: green')
    await client.loginPKCEGrant(gc_clientId, gc_redirectUrl, {})

    //GET Current UserId
    let user = await uapi.getUsersMe({})
    console.log(user)
    userId = user.id
    createWSS(userId)
    //TODO: Enter in additional starting code.
  } catch (err) {
    console.log('Error: ', err)
  }
}

async function createWSS(userId) {
  try {
    //Need to store wss as only can have 15 per agent. Also bad practice to create multiply
    if (sessionStorage.getItem('gc_channelid')) {
      console.log('channelid already exists...')
      var channelid = sessionStorage.getItem('gc_channelid')

      let callsTopic = `v2.users.${userId}.conversations.calls`
      await napi.postNotificationsChannelSubscriptions(channelid, [{ id: callsTopic }])
      console.log(`%cSubscribed to topic ${callsTopic}`, 'color: green')
    } else {
      let channel = await napi.postNotificationsChannels()
      console.log('Created Notification Channel: ', channel)

      let callsTopic = `v2.users.${userId}.conversations.calls`
      await napi.postNotificationsChannelSubscriptions(channel.id, [{ id: callsTopic }])
      console.log(`Subscribed to topic ${callsTopic}`)
      sessionStorage.setItem('gc_channelid', channel.id)
    }
  } catch (err) {
    console.error('Notification Error: ', err)
    console.warn('Clearing channelId')
    sessionStorage.removeItem('gc_channelid')
    createWSS(userId)
  }

  //Create websocket for events
  try {
    let socket = new WebSocket(`wss://streaming.${gc_region}/channels/${sessionStorage.getItem('gc_channelid')}`)
    socket.onmessage = async function(event) {
      let details = JSON.parse(event.data)
      details?.eventBody?.message === 'WebSocket Heartbeat' ? console.log('%c%s Heartbeat', 'color: red', '❤️') : console.log(details)
      //if calls notification
      if (details.topicName.includes('calls')) {
        console.log('Calls Notification: ', details)
        let agentParticipant = details.eventBody.participants.slice().reverse().find(p => p.purpose === 'agent' && p.state === 'connected')
        if (agentParticipant) {
        }
      }
    }
    console.log(`Waiting for events on wss://streaming.${gc_region}/channels/${sessionStorage.getItem('gc_channelid')}`)
  } catch (err) {
    console.error('Websocket error: ', err)
    console.warn('Clearing channelId')
    sessionStorage.removeItem('gc_channelid')
    createWSS(userId)
  }
}
