'use strict' //Enables strict mode is JavaScript
let url = new URL(document.location.href)
let reason = url.searchParams.get('reason')

//USed for voice to digital deflection
if (reason == 'digitalUpgrade') {
  window.history.pushState(null, 'PWA', './home.html')
  sessionStorage.setItem('reason', 'digitalUpgrade')
  setTimeout(function (e) {
    Genesys('command', 'Messenger.open')
  }, 1500)
}

//Used for oauth redirect
if (reason == 'oauth') {
  window.history.pushState(null, 'PWA', './home.html')
  setTimeout(function (e) {
    Genesys('command', 'Messenger.open')
  }, 1500)
  setTimeout(function (f) {
    Genesys('command', 'MessagingService.sendMessage', {
      message: 'logged in',
    })
  }, 2000)
}

//Notification Start----->
Notification.requestPermission().then(function (permission) {
  // If the user accepts, let's create a notification
  if (permission === 'granted') {
    console.log('notification permission granted')
  }
})

//set user login details if any
Genesys('subscribe', 'Messenger.opened', function () {
  console.log('messenger opened')
  try {
    let reason = sessionStorage.getItem('reason')
    let userName = localStorage.getItem('userName')

    //make first char capital
    let firstName = userName.substring(0, userName.indexOf('.'))
    let firstCapital = firstName.charAt(0).toUpperCase() + firstName.slice(1)
    firstName = firstCapital

    //make first char capital
    let lastName = userName.substring(userName.indexOf('.') + 1, userName.indexOf('@'))
    let lastCapital = lastName.charAt(0).toUpperCase() + lastName.slice(1)
    lastName = lastCapital

    Genesys('command', 'Database.set', {
      messaging: {
        customAttributes: {
          firstName: firstName,
          lastName: lastName,
          email: userName,
          reason: reason,
        },
      },
    })
    console.log('userName: ', userName)
  } catch (err) {
    console.error(err)
  }
})

Genesys('subscribe', 'MessagingService.messagesReceived', function (o) {
  //terms and conditions webform
  if (o.data.messages[0].text === 'webform') {
    setTimeout(function () {
      questionForm()
      Genesys('command', 'Messenger.close')
    }, 3000)
    return
  }
  //login form
  if (o.data.messages[0].text === 'You will be directed to login a moment.') {
    setTimeout(function () {
      oauth()
    }, 3000)
    return
  } //You will be directed to login a moment.
  try {
    if (sessionStorage.getItem('gc_widget') === 'false' && o.data.messages[0].direction === 'Outbound') {
      var notification = new Notification('New Message', { tag: 'genesys', body: o.data.messages[0].text, icon: './svgs/anz_icon.png' })
      Genesys('command', 'Messenger.open')
      return
    }
    if (document.hasFocus() === false) {
      var notification = new Notification('New Message', { tag: 'genesys', body: o.data.messages[0].text, icon: './svgs/anz_icon.png' })
      Genesys('command', 'Messenger.open')
    }
  } catch (err) {
    //console.error(err)
  }
})
Genesys('subscribe', 'Conversations.closed', function (o) {
  sessionStorage.setItem('gc_widget', 'false')
})
Genesys('subscribe', 'Conversations.opened', function (o) {
  sessionStorage.setItem('gc_widget', 'true')
})
//Notification End----->

function questionForm() {
  let main = document.getElementById('main')
  let form = `
        <div>
        <fieldset>
          <label for="terms">Terms</label>
          <p>Do you agree with the terms and conditions as per your contract</p>
        </fieldset>
        <fieldset style="display: block">
          <input id="agree" type="checkbox" />
          <label for="checkbox">I Agree</label>
        </fieldset>
        <fieldset style="display: block">
          <input id="disagree" type="checkbox" />
          <label for="checkbox">I Disagree</label>
        </fieldset>
        <div class="formButtons">
        <button onclick="closeForm('Form completed')" class="login">Submit</button>
        </div>
      </div>`

  main.innerHTML = form
}

function closeForm(message) {
  Genesys('command', 'Messenger.open')
  Genesys('command', 'MessagingService.sendMessage', {
    message: message,
  })
  buildHome()
}

function oauth() {
  console.log('oauth')
  document.location.href = './index.html?reason=oauth'
}
