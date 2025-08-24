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
    //TODO: Enter in additional starting code.
  } catch (err) {
    console.log('Error: ', err)
  }
}
