'use strict' //Enables strict mode is JavaScript
let url = new URL(document.location.href)
let reason = url.searchParams.get('reason')

//Notification Start----->
Notification.requestPermission().then(function(permission) {
  // If the user accepts, let's create a notification
  if (permission === 'granted') {
    console.log('notification permission granted')
  }
})

//set user login details if any
Genesys('subscribe', 'Messenger.opened', function() {
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

Genesys('subscribe', 'MessagingService.messagesReceived', function(o) {
  if (o.data.messages[0].text === 'dispute resolved') {
    clearDispute()
    let transaction = sessionStorage.getItem('transaction')
    console.log('removing dispute: ', transaction)
  }
  try {
    if (sessionStorage.getItem('gc_widget') === 'false' && o.data.messages[0].direction === 'Outbound') {
      var notification = new Notification('New Message', { tag: 'genesys', body: o.data.messages[0].text, icon: '.LOGO' })
      Genesys('command', 'Messenger.open')
      return
    }
    if (document.hasFocus() === false) {
      var notification = new Notification('New Message', { tag: 'genesys', body: o.data.messages[0].text, icon: '.LOGO' })
      Genesys('command', 'Messenger.open')
    }
  } catch (err) {
    //console.error(err)
  }
})
Genesys('subscribe', 'Conversations.closed', function(o) {
  sessionStorage.setItem('gc_widget', 'false')
})
Genesys('subscribe', 'Conversations.opened', function(o) {
  sessionStorage.setItem('gc_widget', 'true')
})
//Notification End----->

function closeForm(message) {
  Genesys('command', 'Messenger.open')
  Genesys('command', 'MessagingService.sendMessage', {
    message: message,
  })
  buildHome()
}

