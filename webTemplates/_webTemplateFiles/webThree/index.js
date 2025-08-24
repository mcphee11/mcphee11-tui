'use strict' //Enables strict mode is JavaScript
let url = new URL(document.location.href)
let gc_region = url.searchParams.get('gc_region')
let gc_clientId = url.searchParams.get('gc_clientId')
let gc_redirectUrl = url.searchParams.get('gc_redirectUrl')
let csvRows = []

//Getting and setting the GC details from dynamic URL and session storage
gc_region ? sessionStorage.setItem('gc_region', gc_region) : gc_region = sessionStorage.getItem('gc_region')
gc_clientId ? sessionStorage.setItem('gc_clientId', gc_clientId) : gc_clientId = sessionStorage.getItem('gc_clientId')
gc_redirectUrl ? sessionStorage.setItem('gc_redirectUrl', gc_redirectUrl) : gc_redirectUrl = sessionStorage.getItem('gc_redirectUrl')

let platformClient = require('platformClient')
const client = platformClient.ApiClient.instance
const capi = new platformClient.ConversationsApi()

// Configure Client App
const ClientApp = window.purecloud.apps.ClientApp
const myClientApp = new ClientApp({
  pcEnvironment: gc_region,
})

async function start() {
  try {
    client.setEnvironment(gc_region)
    client.setPersistSettings(true, '_mm_')

    console.log('%cLogging in to Genesys Cloud', 'color: green')
    await client.loginPKCEGrant(gc_clientId, gc_redirectUrl, {})
    getUTCOffset()
    thisMonth()
    // getData() // Uncomment this line to get data on load
  } catch (err) {
    console.log('Error: ', err)
  }
}

function last30days() {
  let today = new Date()
  let aMonthAgo = new Date()
  aMonthAgo.setMonth(aMonthAgo.getMonth() - 1)
  // prettier-ignore
  document.getElementById('datepicker').value = `${aMonthAgo.getFullYear()}-${String(aMonthAgo.getMonth() + 1).padStart(2, '0')}-${String(aMonthAgo.getDate()).padStart(2, '0')}/${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
}

function thisMonth() {
  let today = new Date()
  // prettier-ignore
  document.getElementById('datepicker').value = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-01/${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
}

function getUTCOffset() {
  const now = new Date()
  const offsetMinutes = now.getTimezoneOffset()
  const offsetHours = -offsetMinutes / 60 // Invert for standard UTC offset representation

  let offsetString = ''
  if (offsetHours === 0) {
    offsetString = '+00:00'
  } else {
    const sign = offsetHours > 0 ? '+' : '-'
    const absHours = Math.abs(Math.floor(offsetHours))
    const minutes = Math.abs(Math.floor((offsetHours - absHours) * 60))
    const hoursString = absHours.toString().padStart(2, '0')
    const minutesString = minutes.toString().padStart(2, '0')
    offsetString = `${sign}${hoursString}:${minutesString}`
  }
  document.getElementById('timeZone').value = offsetString
}

// download csv
function buildCsv() {
  let csvContent = 'data:text/csv;charset=utf-8,' + csvRows.map((e) => e.join(',')).join('\r\n')
  const encodedUri = encodeURI(csvContent)
  const link = document.createElement('a')
  link.setAttribute('href', encodedUri)
  link.setAttribute('download', 'report.csv')
  document.body.appendChild(link)
  link.click()
}

async function clearData() {
  console.log('%cClearing data', 'color: green')
  document.getElementById('tableLocation').innerHTML = ''
  document.getElementById('spinner').style.display = 'none'
  document.getElementById('download').style.display = 'none'
  csvRows = []
}

function notification(type, message) {
  if (window.location !== window.parent.location) {
    // if in an iframe
    myClientApp.alerting.showToastPopup(type, message)
    return
  }
  window.alert(message)
  return
}

// dynamic table creation
async function buildTableRows(rows) {
  for (const row of rows) {
    let tableBody = document.getElementById('tableBody')
    if (!tableBody) {
      // create the top table row if its not already there
      console.log('Creating table')
      let top = document.getElementById('tableLocation')
      let guxTable = document.createElement('gux-table')
      let table = document.createElement('table')
      let header = document.createElement('thead')
      let headerRow = document.createElement('tr')
      let tbody = document.createElement('tbody')

      headerRow.setAttribute('data-row-id', 'head')
      tbody.setAttribute('id', 'tableBody')
      table.setAttribute('slot', 'data')
      guxTable.setAttribute('resizable-columns', '')

      header.appendChild(headerRow)
      guxTable.appendChild(table)
      table.appendChild(header)
      table.appendChild(tbody)

      // create column names on first row
      for (const item of row) {
        let th = document.createElement('th')
        th.setAttribute('data-column-name', item)
        th.style.textWrap = 'auto'
        th.innerHTML = item
        th.title = item
        headerRow.appendChild(th)
      }
      top.appendChild(guxTable)
      continue
    }
    if (tableBody) {
      // add data to the row
      let tr = document.createElement('tr')
      for (const item of row) {
        let column = document.createElement('td')
        column.innerHTML = item
        tr.appendChild(column)
      }
      tableBody.appendChild(tr)
    }
  }
}

async function getData() {
  console.log('%cGetting data', 'color: green')
  csvRows = []
  document.getElementById('tableLocation').innerHTML = ''
  document.getElementById('spinner').style.display = 'block'
  // TODO: ENTER IN YOUR CODE HERE (EXAMPLE WITH PAGINATION)
  let pageNumber = 1
  try {
    const conversations = await getConversations(pageNumber)
    console.log('Conversations page1: ', conversations)
    if (conversations.totalHits > 100) {
      while (pageNumber < Math.ceil(conversations.totalHits / 100)) {
        pageNumber++
        const nextConversations = await getConversations(pageNumber)
        conversations.conversations = conversations.conversations.concat(nextConversations.conversations)
      }
    }

    for (const conv of conversations.conversations) {
      // add the column names to the first row
      if (csvRows.length === 0) {
        csvRows.push(['ConversationId'])
        csvRows[0].push('Media Type')
        csvRows[0].push('Start Date')
      }
      // add the data to the row
      csvRows.push([conv.conversationId])
      csvRows[csvRows.length - 1].push(conv.participants[0].sessions[0].mediaType)
      csvRows[csvRows.length - 1].push(new Date(conv.conversationStart).toLocaleString().replace(',', ' '))
    }
    console.log(csvRows)
    await buildTableRows(csvRows)
    document.getElementById('spinner').style.display = 'none'
    document.getElementById('download').style.display = 'block'
  } catch (err) {
    console.log('Error: ', err)
    notification('Error', `Error: ${err}`)
  }
}

async function getConversations(pageNumber) {
  // TODO: your query
  // ------------- EXAMPLE ------------------
  let conversations = await capi.postAnalyticsConversationsDetailsQuery({
    // prettier-ignore
    interval: `${document.getElementById('datepicker').value.split('/')[0]}T00:00:00${document.getElementById('timeZone').value}/${document.getElementById('datepicker').value.split('/')[1]}T23:59:59${document.getElementById('timeZone').value}`,
    paging: {
      pageSize: 100,
      pageNumber: pageNumber,
    },
    segmentFilters: [
      {
        type: 'and',
        predicates: [
          {
            dimension: 'mediaType',
            value: 'voice',
          },
        ],
      },
    ],
    orderBy: 'conversationStart',
  })
  // ------------- EXAMPLE ------------------
  return conversations
}
