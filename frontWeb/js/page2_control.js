var model2=angular.module('MyApp.page2', ['ui.bootstrap']);
model2.controller('p2mainCtrl',[
	'$scope',
	'$modal',
	'$rootScope',
	function ($scope,$rootScope,$modal){

	}
]);



model2.controller('text2',[
	'$scope',
	'$modal',//注入modal
	'$log',
	function ($scope,$modal,$log){

		$scope.aa="30";

		//单击卡片时执行的东西,index代表设备列表的第几个
		$scope.showdetail=function(size,myindex){
			var modalInstance = $modal.open({
      			templateUrl: 'myModalContent.html',
      			controller: 'ModalInstanceCtrl',
      			size: size,
      			myindex:myindex,
      			//backdrop:"static",//点击空白处不会退出modal
      			//keyboard :false, //按键盘esc不会退出modal
      			resolve: {//向模态框发送一些值，通过注入
        			items: function () {//一个属性或对象可以在modal的控制器里面注入
          				return $scope.items;
        			}
      			}
    		});
			console.log(modalInstance.myindex);
			//处理模态框返回的结果，selectedItem从modal中返回
			//回调结束后，then，这里的回调就是出现的modal
    		modalInstance.result.then(function (selectedItem) {
      			$scope.selected = selectedItem;//如果返回成功
    		}, function () {
      			$log.info('Modal dismissed at: ' + new Date());
    		});
		}


		$scope.items = ['item1', 'item2', 'item3'];
	}
]);

model2.controller('ModalInstanceCtrl',[
	'$scope',
	'$modalInstance',//父控制器的注入，可以提供close与dismiss方法
	'items',//这个items是父控制器resolve返回并可以注入的
	function ($scope, $modalInstance, items) {

		$scope.items = items;
		$scope.myindex = $modalInstance.myindex;
		$scope.selected = {//会在前端的ngclick中改变
			item: $scope.items[0]//会产生scope.selected.item的值，赋值默认为items[0]
		};

		$scope.ok = function () {
			//向调用modal的控制器返回一些东西
			$modalInstance.close($scope.selected.item);
		};

  		$scope.cancel = function () {
    		$modalInstance.dismiss('cancel');
  		};
	}
]);