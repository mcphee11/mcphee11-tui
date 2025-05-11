'use strict' //Enables strict mode is JavaScript
let userName = localStorage.getItem('userName')
let transaction = null

// check for login user
if (document.location.href.includes('index.html')) {
  userName ? (document.location.href = './home.html' + document.location.search) : null
}

if (document.location.href.includes('home.html')) {
  let userName = localStorage.getItem('userName')
  let firstName
  if (userName == 'null' || userName == '') {
    firstName = 'Guest'
  } else {
    firstName = userName.substring(0, userName.indexOf('.'))
    let capital = firstName.charAt(0).toUpperCase() + firstName.slice(1)
    firstName = capital
  }
  let name = document.getElementById('name')
  name.innerText = firstName
  buildHome()
}

//Request Notification permission
Notification.requestPermission().then(function (permission) {
  if (permission === 'granted') {
    console.log('notification permission granted')
  }
})

//Check browser support for SW
if ('serviceWorker' in navigator) {
  navigator.serviceWorker
    .register('service-worker.js')
    .then(function (registration) {
      console.info('Service Worker Registered.')
    })
    .catch(function (err) {
      console.error(err)
    })
}

if ('PushManager' in window) {
  console.info('push supported')
}

function serviceWorkerMessage(message) {
  console.log(message)
}

navigator.serviceWorker.addEventListener('message', async (event) => {
  // Optional: ensure the message came from workbox-broadcast-update
  if (event.data.meta === 'workbox-broadcast-update') {
    console.log('update received', event)
    if (event.data.payload.updatedURL === document.location.href) {
      document.location.reload()
    }
  }
  if (event.data.badge) {
    console.log(event.data)
  }
})

document.addEventListener('click', (e) => {
  //logout clicked
  if (e.target.id.startsWith('logout')) {
    document.title = 'Logout'
    document.getElementById('account_img').className = 'icon_bottom'
    document.getElementById('doc_img').className = 'icon_bottom'
    document.getElementById('home_img').className = 'icon_bottom'
    document.getElementById('profile_img').className = 'icon_bottom'
    document.getElementById('main').innerHTML = ''
    localStorage.removeItem('userName')
    localStorage.removeItem('userTelNumber')
    sessionStorage.removeItem('reason')
    document.location.href = './index.html'
  }
  //support clicked
  if (e.target.id.startsWith('support')) {
    document.title = 'Support'
    Genesys('command', 'Messenger.open')
    document.getElementById('account_img').className = 'icon_bottom'
    document.getElementById('doc_img').className = 'icon_bottom'
    document.getElementById('home_img').className = 'icon_bottom'
    document.getElementById('profile_img').className = 'icon_bottom'
  }
  //account clicked
  if (e.target.id.startsWith('account')) {
    document.title = 'Accounts'
    document.getElementById('account_img').className = 'icon_bottom_blue'
    document.getElementById('doc_img').className = 'icon_bottom'
    document.getElementById('home_img').className = 'icon_bottom'
    document.getElementById('profile_img').className = 'icon_bottom'
    document.getElementById('main').innerHTML = ''
    buildAccounts()
  }
  //doc clicked
  if (e.target.id.startsWith('doc')) {
    document.title = 'Transactions'
    document.getElementById('account_img').className = 'icon_bottom'
    document.getElementById('doc_img').className = 'icon_bottom_blue'
    document.getElementById('home_img').className = 'icon_bottom'
    document.getElementById('profile_img').className = 'icon_bottom'
    buildTransactions()
  }
  //anz plus clicked
  if (e.target.id.startsWith('home_img')) {
    document.title = 'ANZ Plus'
    document.getElementById('account_img').className = 'icon_bottom'
    document.getElementById('doc_img').className = 'icon_bottom'
    document.getElementById('home_img').className = 'icon_bottom_blue'
    document.getElementById('profile_img').className = 'icon_bottom'
    buildHome()
  }
  //delete clicked
  if (e.target.id.startsWith('delete')) {
    //2d25a95b-be36-4ee6-b1f4-df3da96b1c41
    localStorage.removeItem(`_2d25a95b-be36-4ee6-b1f4-df3da96b1c41:gcmcopn`)
    localStorage.removeItem(`_2d25a95b-be36-4ee6-b1f4-df3da96b1c41:gcmcsessionActive`)
    localStorage.removeItem('userName')
    localStorage.removeItem('userTelNumber')
    sessionStorage.removeItem('reason')
    Genesys('command', 'Identifiers.purgeAll', {})
    document.location.href = './index.html'
  }
  //profile is clicked
  if (e.target.id.startsWith('profile')) {
    document.title = 'Profile'
    document.getElementById('account_img').className = 'icon_bottom'
    document.getElementById('doc_img').className = 'icon_bottom'
    document.getElementById('home_img').className = 'icon_bottom'
    document.getElementById('profile_img').className = 'icon_bottom_blue'
    document.getElementById('main').innerHTML = ''
    buildProfile()
  }
  //modal close button is clicked for either modal
  if (e.target.id.startsWith('modalClose')) {
    document.getElementById('modal_dispute').close()
    document.getElementById('account_img').className = 'icon_bottom'
    document.getElementById('doc_img').className = 'icon_bottom'
    document.getElementById('home_img').className = 'icon_bottom'
    document.getElementById('profile_img').className = 'icon_bottom'
  }

  // more button is clicked on card
  if (e.target.id.startsWith('button_')) {
    let element = e.target.id.substring(7, e.target.id.length)
    if (document.getElementById(`drop_${element}`).style.display == 'block') {
      document.getElementById(`drop_${element}`).style.display = 'none'
    } else {
      document.getElementById(`drop_${element}`).style.display = 'block'
    }
  }

  // dispute button is clicked on more drop down
  if (e.target.id.startsWith('dropdown_dispute_')) {
    document.title = 'Dispute'
    let element = e.target.id.substring(17, e.target.id.length)
    document.getElementById(`drop_${element}`).style.display = 'none'
    document.getElementById('modal_dispute').showModal()
    transaction = element
  }

  //modal phone button clicked
  if (e.target.id.startsWith('phone')) {
    sessionStorage.setItem('reason', document.getElementById(transaction).children[1].children[1].innerText)
    document.getElementById('modal_dispute').close()
    console.log(transaction)
    Genesys('command', 'Journey.record', { eventName: 'Message_Dispute', customAttributes: { Transaction: `${document.getElementById(transaction).children[1].children[1].innerText}` } })
    clickToCallAuth()
  }

  // modal chat button clicked
  if (e.target.id.startsWith('message')) {
    sessionStorage.setItem('reason', document.getElementById(transaction).children[1].children[1].innerText)
    document.getElementById('modal_dispute').close()
    console.log(transaction)
    Genesys('command', 'Journey.record', { eventName: 'Message_Dispute', customAttributes: { Transaction: `${document.getElementById(transaction).children[1].children[1].innerText}` } })
    Genesys(
      'command',
      'Messenger.open',
      {},
      () => {
        setTimeout(function (f) {
          Genesys('command', 'MessagingService.sendMessage', {
            message: `I'm enquiring about the transaction: ${document.getElementById(transaction).children[1].children[1].innerText}`,
          })
        }, 2000)
      },
      (error) => {
        console.log("Couldn't open messenger.", error)
      }
    )
  }
  // modal call button was clicked
  if (e.target.id.startsWith('phone')) {
    document.getElementById('modal_dispute').close()
    Genesys('command', 'Journey.record', { eventName: 'Phone_Call', customAttributes: { Transaction: `${document.getElementById(transaction).children[1].children[1].innerText}` } })
  }

  // ---- KEY PAD -----
  if (e.target.id === 'keypad_button') {
    let pin = document.getElementById('pin')
    let update = pin.innerText.replace('_', e.target.innerText)
    pin.innerText = update
    if (!document.getElementById('pin').innerText.includes('_')) {
      console.log('4 Digits Entered')
      console.log(document.getElementById('pin').innerText)
      if (document.getElementById('pin').innerText == '1 2 3 4') {
        console.log('login')
        document.getElementById('pin').innerText = '_ _ _ _'
        localStorage.setItem('userName', 'USER_NAME')
        localStorage.setItem('userTelNumber', 'tel:USER_NUMBER')
        document.location.href = './home.html' + document.location.search
      }
      // Add your own additional userName & userTelNumber below for additional user logins
      // if (document.getElementById('pin').innerText == '4 3 2 1') {
      //   console.log('login')
      //   document.getElementById('pin').innerText = '_ _ _ _'
      //   localStorage.setItem('userName', 'demo@test.com')  // add your email here
      //   localStorage.setItem('userTelNumber', 'tel:+61400000000')  // add your ani here
      //   document.location.href = './home.html' + document.location.search
      // }
      document.getElementById('pin').innerText = '_ _ _ _'
    }
  }
  if (e.target.id.includes('backspace')) {
    console.log('backspace')
    let pin = document.getElementById('pin')
    pin.innerText = '_ _ _ _'
  }
})

//sends message to service worker returns promise
async function sendMessage(message) {
  return new Promise(function (resolve, reject) {
    var messageChannel = new MessageChannel()
    messageChannel.port1.onmessage = function (event) {
      if (event.data.error) {
        reject(event.data.error)
      } else {
        resolve(event.data)
        serviceWorkerMessage(event.data)
      }
    }
    navigator.serviceWorker.controller.postMessage(message, [messageChannel.port2])
  })
}

async function clickToCallAuth() {
  let body = {
    KEY: localStorage.getItem('userTelNumber'),
    timeString: 'doneInServer',
    reason: sessionStorage.getItem('reason'),
    email: localStorage.getItem('userName'),
  }

  let tableUpdate = await fetch('UPDATE_DATATABLE_CLOUD_FUNCTION_URL', {
    method: 'POST',
    headers: {
      'content-Type': 'application/json',
      Authorization: '', // TODO add Basic Auth
    },
    body: JSON.stringify(body),
  })

  let response = await tableUpdate.json()
  console.log(response)
}

function buildHome() {
  let main = document.getElementById('main')
  let html = `<div style="width: 100%; display: flex">
  <img id="safety_img" src="./LOGO" style="width: 200px; margin: auto;" />
  </div>`
  main.innerHTML = ''
  main.innerHTML = html
}

function buildAccounts() {
  let main = document.getElementById('main')
  let html = `<div style="width: 100%">
  <div id="card1" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/piggy.webp" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Savings</h3>
        <p class="account_middle_text"><strong>+$123,050.55</strong> Balance</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card1" class="account_button">
        <img id="button_card1" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card1" class="dropdown-content">
          <a href="#Details">Details</a>
      </div>
      </div>
  </div>

    <div id="card2" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/card.jpg" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Credit Card</h3>
        <p class="account_middle_text"><strong>-$4,281.30</strong> Balance</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card2" class="account_button">
        <img id="button_card2" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card2" class="dropdown-content">
          <a href="#Details">Details</a>
      </div>
      </div>
  </div>

  </div>`
  main.innerHTML = ''
  main.innerHTML = html
}

function buildProfile() {
  let main = document.getElementById('main')
  let html = `<div style="text-align: center; width: 90%; margin-top: 20px">
      <p>Built as a POC by Genesys</p>
      <div style="display: flex; justify-content: space-evenly">
      <div>
        <a id="delete" href="#delete" style="text-decoration: none;">
          <img id="delete_img" src="./svgs/delete_black.svg" style="margin: 10px; width: 50px" />
        </a>
        <p>Delete session</p>
      </div>
      <div>
        <a id="logout" href="#logout">
          <img id="logout_img" src="./svgs/login_24dp.svg" style="margin: 10px; width: 50px" />
        </a>
        <p>Logout session</p>
      </div>
      </div>
      <p style="font-size: xx-small; margin-top: 20px">version 2.3</p>
    </div>`
  main.innerHTML = ''
  main.innerHTML = html
}

function buildTransactions() {
  let main = document.getElementById('main')
  let html = `<div style="width: 100%">
  <div id="card1" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/card.jpg" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Credit Card</h3>
        <p class="account_middle_text"><strong>-$103.95</strong> Woolworths</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card1" class="account_button">
        <img id="button_card1" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card1" class="dropdown-content">
          <a href="#Details">Details</a>
          <a href="#Dispute" id="dropdown_dispute_card1">Dispute</a>
      </div>
      </div>
  </div>

    <div id="card2" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/card.jpg" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Credit Card</h3>
        <p class="account_middle_text"><strong>-$81.20</strong> Ebay</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card2" class="account_button">
        <img id="button_card2" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card2" class="dropdown-content">
          <a href="#Details">Details</a>
          <a href="#Dispute" id="dropdown_dispute_card2">Dispute</a>
      </div>
      </div>
  </div>

    <div id="card3" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/card.jpg" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Credit Card</h3>
        <p class="account_middle_text"><strong>-$5.40</strong> Cafe Now</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card3" class="account_button">
        <img id="button_card3" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card3" class="dropdown-content">
          <a href="#Details">Details</a>
          <a href="#Dispute" id="dropdown_dispute_card3">Dispute</a>
      </div>
      </div>
  </div>

    <div id="card4" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/piggy.webp" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Savings</h3>
        <p class="account_middle_text"><strong>+$300.00</strong> Thanks for dinner</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card4" class="account_button">
        <img id="button_card4" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card4" class="dropdown-content">
          <a href="#Details">Details</a>
          <a href="#Dispute" id="dropdown_dispute_card4">Dispute</a>
      </div>
      </div>
  </div>

    <div id="card5" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/card.jpg" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Credit Card</h3>
        <p class="account_middle_text"><strong>-$57.40</strong> JB-HiFi</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card5" class="account_button">
        <img id="button_card5" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card5" class="dropdown-content">
          <a href="#Details">Details</a>
          <a href="#Dispute" id="dropdown_dispute_card5">Dispute</a>
      </div>
      </div>
  </div>

    <div id="card6" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/piggy.webp" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Savings</h3>
        <p class="account_middle_text"><strong>+$1,100.00</strong> Transfer</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card6" class="account_button">
        <img id="button_card6" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card6" class="dropdown-content">
          <a href="#Details">Details</a>
          <a href="#Dispute" id="dropdown_dispute_card6">Dispute</a>
      </div>
      </div>
  </div>

    <div id="card6" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/card.jpg" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Credit Card</h3>
        <p class="account_middle_text"><strong>-$11.00</strong> Cafe Now</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card6" class="account_button">
        <img id="button_card6" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card6" class="dropdown-content">
          <a href="#Details">Details</a>
          <a href="#Dispute" id="dropdown_dispute_card6">Dispute</a>
      </div>
      </div>
  </div>

    <div id="card7" class="account_row">
        <div id="left" class="account_left">
        <img src="./svgs/card.jpg" style="width: 75px" />
      </div>
      <div id="center" class="account_middle">
        <h3 class="account_middle_text">Credit Card</h3>
        <p class="account_middle_text"><strong>-$230.00</strong> Woolworths</p>
      </div>
      <div id="right" class="account_right dropdown">
        <button id="button_card7" class="account_button">
        <img id="button_card7" src="./svgs/more_white.svg" style="width: 40px; filter: invert(36%) sepia(97%) saturate(3855%) hue-rotate(181deg) brightness(91%) contrast(95%);" />
        </button>
          <div id="drop_card7" class="dropdown-content">
          <a href="#Details">Details</a>
          <a href="#Dispute" id="dropdown_dispute_card7">Dispute</a>
      </div>
      </div>
  </div>

  </div>`
  main.innerHTML = ''
  main.innerHTML = html
}
