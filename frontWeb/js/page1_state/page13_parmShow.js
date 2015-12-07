var page13_model=angular.module('MyApp.page13',["xeditable","checklist-model"]);

page13_model.run(function(editableOptions) { //插件xeditable的初始化样式
  editableOptions.theme = 'bs3'; // bootstrap3 theme. Can be also 'bs2', 'default'
});



/*================================通风等级表============================================
	作用：1.当http请求数据后，马上更新以checklist显示的风机组模型，同时也作为显示在表格的模型
		  2.uploadData()上传当前的data数据，保存参数设置
==============================================================================*/
page13_model.controller('page13_wenduTableCtrl', [
	'$scope',
	'$http',
	'$filter',
	function ($scope,$http,$filter){
		//风机选择后的状态,因为有20个等级所以有20个组，会在得到json的get请求中进行同步
		$scope.fenjistat = [
			[],[],[],[],[],
			[],[],[],[],[],
			[],[],[],[],[],
			[],[],[],[],[]
		];
		//$http.get('/testjson/page13/WindTables.json').success(function(data){ //测试用json
		$http.get('/resetful/nm820/sysPara/WindTables').success(function(data){
			$scope.data=data;

			//同步风机组checklist,20个通风等级，循环0-19
			for (var i = 0; i <= 19; i++) {
				var fan=$scope.data.WindTables[i].Fan;
				$scope.fenjistat[i]=[];//显示的checklist全部不勾选,按照所勾选的风机等级
	  			//如果低位有值 如100111&000001 =0000001 //返回的是1号风机在，那么显示checklist时也要打勾
	  			if((fan&1)==1){$scope.fenjistat[i].push(1);};
	  			if((fan&2)==2){$scope.fenjistat[i].push(2);};
	  			if((fan&4)==4){$scope.fenjistat[i].push(3);};
	  			if((fan&8)==8){$scope.fenjistat[i].push(4);};
	  			if((fan&16)==16){$scope.fenjistat[i].push(5);};
	  			if((fan&32)==32){$scope.fenjistat[i].push(6);};
	  			if((fan&64)==64){$scope.fenjistat[i].push(7);};
	  			if((fan&128)==128){$scope.fenjistat[i].push(8);};
			};
		});
		//风机选择checklist的模型	
		$scope.fenjistats = [
    		{value: 1, text: '1'},
    		{value: 2, text: '2'},
    		{value: 3, text: '3'},
    		{value: 4, text: '4'},
    		{value: 5, text: '5'},
    		{value: 6, text: '6'},
    		{value: 7, text: '7'},
    		{value: 8, text: '8'}
  		];

  		//用xeditable修改完Fan的值后，每个通风等级的风机，反向更新会data[i].Fan
  		//将checklist模型转化为更新原来的原始fan二进制表示，$scope.fenjistat变回data[i].Fan
  		$scope.onafterSetFan=function(index){
  			$scope.data.WindTables[index].Fan=0;//先清空
  			for (var j = 0; j < $scope.fenjistat[index].length; j++) {//每个checklist进行循环，j代表每个通风等级里面的风机索引
  		 		if($scope.fenjistat[index][j]==1){$scope.data.WindTables[index].Fan+=1;}
  		 		if($scope.fenjistat[index][j]==2){$scope.data.WindTables[index].Fan+=2;}
  		 		if($scope.fenjistat[index][j]==3){$scope.data.WindTables[index].Fan+=4;}
  		 		if($scope.fenjistat[index][j]==4){$scope.data.WindTables[index].Fan+=8;}
  		 		if($scope.fenjistat[index][j]==5){$scope.data.WindTables[index].Fan+=16;}
  		 		if($scope.fenjistat[index][j]==6){$scope.data.WindTables[index].Fan+=32;}
  		 		if($scope.fenjistat[index][j]==7){$scope.data.WindTables[index].Fan+= 64;}
  		 		if($scope.fenjistat[index][j]==8){$scope.data.WindTables[index].Fan+= 128;}
  		 	};
  		 	//console.log($scope.data[index].Fan);
  		};

  		//====================保存修改后的参数设置,但要注意经过xeditable后数值是字符型的不是数值型的================
  		$scope.uploadData=function(){
  			$http.post('/resetful/nm820/sysPara/WindTables',$scope.data)
  			.success(function(){
  				alert("保存修改成功");
  			})
  		}

	}
]);

/*============================控制器 温度曲线======================================
	1.图表切换
	2.请求温度曲线数据
=======================================================================================*/
page13_model.controller('page13_wenduCurveCtrl', [
	'$scope',
	'$http',
	function ($scope,$http){
		$scope.ifShowTab=true;//一开始时是显示表格不显示曲线
		$scope.ifShowCurve=false;

		//图表显示数据
		$scope.CurveData={
			"series": ["目标温度","加热温度","制冷温度"],
        	//也可以用上面的scope.data的形式
        	"data": [[1,2,3,45,6,7,7,8,8],
        			 [2,34,5,5,5,5,5,5,5],
        			 [3,4,5,5,66,77,88,77]],
        	"labels":[1,2,3,4,5,6,7,8,9],
		};

		$http.get('/resetful/nm820/sysPara/WenduCurve').success(function(data){
			$scope.data=data;
			//温度曲线图表的样子
			$scope.CurveData.labels=data.Day;
			$scope.CurveData.data[0]=data.Target;
			$scope.CurveData.data[1]=data.Heat;
			$scope.CurveData.data[2]=data.Cool;
		});


		//点击显示表格与图表切换
		$scope.changeTabChat = function(){
	    	$scope.ifShowTab=!$scope.ifShowTab;
	    	$scope.ifShowCurve=!$scope.ifShowCurve;
		};

		 //====================保存修改后的参数设置,但要注意经过xeditable后数值是字符型的不是数值型的================
  		$scope.uploadData=function(){
  			$http.post('/resetful/nm820/sysPara/WenduCurve',$scope.data)
  			.success(function(){
  				alert("保存修改成功");
  			})
  		}

  		$scope.string2int=function(){

  		}
	}
]);
/*===========================控制器 最大最小通风等级===============================================
	第三个表格通风等级曲线的控制器，
============================================================================*/
page13_model.controller('page13_WindLevelCtrl', [
	'$scope',
	'$http',
	function ($scope,$http){
		$scope.ifShowTab=true;//一开始时是显示表格不显示曲线
		$scope.ifShowCurve=false;

		//图表显示数据
		$scope.CurveData={
			"series": ["最小通风等级","最大通风等级"],
        	//也可以用上面的scope.data的形式
        	"data": [[1,2,3,45,6,7,7,8,8],
        			 [2,34,5,5,5,5,5,5,5]],
        	"labels":[1,2,3,4,5,6,7,8,9],
		};

		$http.get('/resetful/nm820/sysPara/WindLevel').success(function(data){
			$scope.data=data;

			//温度曲线图表的样子
			$scope.CurveData.labels=data.Day;
			$scope.CurveData.data[0]=data.Min;
			$scope.CurveData.data[1]=data.Max;
		});


		//点击显示表格与图表切换
		$scope.changeTabChat = function(){
	    	$scope.ifShowTab=!$scope.ifShowTab;
	    	$scope.ifShowCurve=!$scope.ifShowCurve;
		};

		//====================保存修改后的参数设置,但要注意经过xeditable后数值是字符型的不是数值型的================
  		$scope.uploadData=function(){
  			$http.post('/resetful/nm820/sysPara/WindLevel',$scope.data)
  			.success(function(){
  				alert("保存修改成功");
  			})
  		}

		$scope.test=function(){
			 return "Username should be `awesome` or `error`";
		};
	}
]);

/*====================系统参数总集成，请注意这里的cotrol将会被多个表使用=================================
  参数：
    1.relayList 用来继电器编码的下拉选项
  方法：
      1.reflashData更新温度月数据
=========================================================*/
page13_model.controller('page13_SysvalCtrl', [
  '$scope',
  '$http',
  '$filter',
  function ($scope,$http,$filter){
    $http.get('/resetful/nm820/sysPara/SysValTable').success(function(data){
      $scope.sysVal=data;
    });

    //修改后用post上传修改
    $scope.uploadData=function(){
      console.log()
      $http.post('/resetful/nm820/sysPara/SysValTable',$scope.sysVal)
        .success(function(){
          alert("保存修改成功");
      })
    };

    //===============下面的代码在继电器编码的选项中用到 表4=======================
    $scope.relayList = [
      {value: 0, text: '0--手动'},
      {value: 1, text: '1--风机组1'},
      {value: 2, text: '2--风机组2'},
      {value: 3, text: '3--风机组3'},
      {value: 4, text: '4--风机组4'},
      {value: 5, text: '5--风机组5'},
      {value: 6, text: '6--风机组6'},
      {value: 7, text: '7--风机组7'},
      {value: 8, text: '8--风机组8'},
      {value: 9, text: '9--加热'},
      {value: 10, text: '10--冷却水泵'},
      {value: 11, text: '11--喷雾'},
      {value: 12, text: '12--回流阀'},
      {value: 13, text: '13--照明1'},
      {value: 14, text: '14--照明2'},
      {value: 15, text: '15--侧风窗开'},
      {value: 16, text: '16--侧风窗关'},
      {value: 17, text: '17--幕帘开'},
      {value: 18, text: '18--幕帘关'},
      {value: 19, text: '19--卷帘1开'},
      {value: 20, text: '20--卷帘1关'},
      {value: 21, text: '21--卷帘2开'},
      {value: 22, text: '22--卷帘2关'},
      {value: 23, text: '23--卷帘3开'},
      {value: 24, text: '24--卷帘3关'},
      {value: 25, text: '25--卷帘4开'},
      {value: 26, text: '26--卷帘4关'},
      {value: 27, text: '27--喂料'},
      {value: 28, text: '28--额外系统1'},
      {value: 29, text: '29--额外系统2'},
      {value: 30, text: '30--额外系统3'},
      {value: 31, text: '31--额外系统4'}
    ]; 

    $scope.relayShow = function(myvalue) {
      var selected = $filter('filter')($scope.relayList, {value: myvalue});//过滤器在relaylist中找出value=myvalue的对象
      return (selected.length) ? selected[0].text : 'Not set';
    };
    //====================================================================

    //==============================下面的代码在控制模式的选项中用到 表7====================================
    $scope.modeList=[
      {value: 0, text: '0--自设'},
      {value: 1, text: '1--曲线'},
      {value: 2, text: '2--曲线加通风'}
    ];

    $scope.modeShow = function(myvalue) {
      var selected = $filter('filter')($scope.modeList, {value: myvalue});//过滤器在relaylist中找出value=myvalue的对象
      return (selected.length) ? selected[0].text : 'Not set';
    };
    //==================================================================================
  }
]);