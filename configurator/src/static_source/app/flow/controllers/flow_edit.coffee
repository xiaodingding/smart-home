angular
.module('appControllers')
.controller 'flowEditCtrl', ['$scope', 'Message', '$stateParams', 'Flow', '$state', 'Workflow', '$timeout'
'log', 'Notify', 'Worker', '$http'
($scope, Message, $stateParams, Flow, $state, Workflow, $timeout, log, Notify, Worker, $http) ->
  vm = this

  # vars
  $scope.callback = {}
  $scope.workflows = []
  $scope.flow = {}
  $scope.elementScripts = {}

  # watcher
  #------------------------------------------------------------------------------
  instance = $scope.$on '$stateChangeStart', (event, toState, toParams, fromState, fromParams, options)->
    if !confirm('Вы точно хотите покинут редактирование процесса?')
      event.preventDefault()

  #------------------------------------------------------------------------------
  # workflow list
  #------------------------------------------------------------------------------
  getWorkflow =->
    success = (result)->
      $scope.workflows = result.items
    error = (result)->
      Message result.data.status, result.data.message
    Workflow.all {}, success, error

  #------------------------------------------------------------------------------
  # flow
  #------------------------------------------------------------------------------
  getFlow =->
    success = (flow) ->
      $scope.flow = flow
      if !$scope.flow?.workers
        $scope.flow.workers = []

      # scripts
      angular.forEach $scope.flow.objects, (object)->
        $scope.elementScripts[object.id] = object.script if object.script?.id?

      $timeout ()->
        $scope.getStatus().then (result)->
          angular.forEach $scope.flows, (value, id)->
            if flow.id == parseInt(id, 10)
              $scope.flow.state = value
      , 500

    error = ->
      $state.go 'dashboard.flow.index'

    Flow.get_redactor {id: $stateParams.id}, success, error
    $scope.remove =->
      if confirm('точно удалить процесс?')
        remove()

  #------------------------------------------------------------------------------
  # remove
  #------------------------------------------------------------------------------
  remove =->
    instance()
    success =->
      $state.go("dashboard.flow.index")
    error =(result)->
      Message result.data.status, result.data.message
    $scope.flow.$delete success, error

  #------------------------------------------------------------------------------
  # save
  #------------------------------------------------------------------------------
  $scope.submit =->
    success =(data)->
      instance()
      Notify 'success', 'Схема успешно сохранена', 3

    error =(result)->
      Message result.data.status, result.data.message

    scheme = $scope.callback.save()
    $scope.flow.objects = scheme.objects || []

    # scripts
    angular.forEach $scope.flow.objects, (object)->
      object.script = $scope.elementScripts[object.id] || null

    $scope.flow.connectors = scheme.connectors || []
    Flow.update_redactor {id: $stateParams.id}, $scope.flow, success, error

  #------------------------------------------------------------------------------
  # init
  #------------------------------------------------------------------------------
  getWorkflow()
  getFlow()

  vm
]