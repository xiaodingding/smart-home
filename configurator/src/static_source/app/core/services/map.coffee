angular
.module('appServices')
.factory 'Map', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/map/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data?.map || data

    create:
      method: 'POST'
      responseType: 'json'
      transformResponse: (data) ->
        data?.map || data

    update:
        method: 'PUT'
        responseType: 'json'

    delete:
      method: 'DELETE'

    all:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        meta: data?.meta || {}
        items: data?.maps || []
]
