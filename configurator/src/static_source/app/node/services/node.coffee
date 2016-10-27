angular
.module('appServices')
.factory 'Node', ['$resource', ($resource) ->
  $resource window.server_url + '/api/v1/node/:id', {id: '@id'},
    show:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        data.node

    create:
      method: 'POST'
      responseType: 'json'

    update:
        method: 'PUT'

    delete:
      method: 'DELETE'

    all:
      method: 'GET'
      responseType: 'json'
      transformResponse: (data) ->
        meta: data.meta
        items: data.nodes

]