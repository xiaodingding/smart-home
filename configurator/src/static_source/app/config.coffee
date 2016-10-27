angular
.module('app')
.config ['$translatePartialLoaderProvider', '$translateProvider'
($translatePartialLoaderProvider, $translateProvider) ->

  window.server_url = 'http://127.0.0.1:3000'

  $translatePartialLoaderProvider.addPart('dashboard');

  $translateProvider.useLoader('$translatePartialLoader', {
    urlTemplate: '/static/translates/{part}/{lang}.json'
    loadFailureHandler: 'LocaleErrorHandler'
  })

  $translateProvider.preferredLanguage 'ru'
  $translateProvider.useSanitizeValueStrategy null
]