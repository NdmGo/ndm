function showMsg(msg, callback ,icon, time){
	if (typeof time == 'undefined'){
		time = 2000;
	}

	if (typeof icon == 'undefined'){
		icon = {};
	}
	var loadT = layer.msg(msg, icon);
	setTimeout(function() {
		layer.close(loadT);
		if (typeof callback == 'function'){
			callback();
		}
	}, time);
}

//字符串转数组对象
function string2ArrayObject(str){
    var data = {};
    kv = str.split('&');
    for(i in kv){
        v = kv[i].split('=');
        data[v[0]] = v[1];
    }
    return data;
}

//表单多维转一维
function array2arr(sa){
    var t = {}

    for (var i = 0; i < sa.length; i++) {
        t[sa[i]['name']] = sa[i]['value'];
    }
    return t;
}

layui.use(['layer','form','element','jquery','table','laydate','util'],function() {
///
var $ = layui.$;
var layer = layui.layer;
var form = layui.form;
var element = layui.element;
var table = layui.table;
var laydate = layui.laydate;
var util = layui.util;

var laytpl = layui.laytpl;
laytpl.config({open: '{@', close: '@}'});

// 设置请求默认值
$.ajaxSetup({
    beforeSend: function (xhr) {
        // 将token塞进Header里
        xhr.setRequestHeader('Authorization', 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9');
        xhr.setRequestHeader('Powered-By', 'NDM');
    },
    complete: function (xhr) {
        // 设置登陆拦截
        // if (xhr.responseJSON.code == "error_unauth") {
        //     console.log("没有登录！");
        //     layer.msg("没有登录！");
        // } else {
        //     console.log("已经登录！");
        // }
    },
});


//监听table表单搜索
form.on('submit(table-sreach)', function (data) {
	var _id = $(this).data('id');

    console.log(data);
    table.reload(_id,{where: data.field, page:{curr: 1}});
});


//监听switch
form.on('switch(*)', function(data){
    var _zt = data.elem.checked ? 'yes' : 'no';
    var _link = $(this).data('link');
    var _id = $(this).data('id');
    var index = layer.load();
    $.post(_link, {'zt':_zt,'id':_id}, function(res) {
        layer.close(index);
        if(res.code == 1){
            layer.msg(res.msg,{icon: 1});
        }else{
            var flag = $("input[name='switch']").prop("checked");
            $("input[name='switch']").prop("checked",!flag);
            form.render("checkbox");
            layer.msg(res.msg,{icon: 2,shift:6});
        }
    },'json');
    return false;
});
    

//监听全局表单提交
form.on('submit(submit_save)', function(data){
	var index = layer.load();
	$.ajax({
        type: "POST",
        headers: {'Content-Type': 'application/json'},
        url: data.form.action,
        data: JSON.stringify(data.field),
        dataType: 'json',
        success: function(res) {
        	layer.close(index);
            if (res.code != 200 ){
                layer.msg(res.message, {icon: 2});
            } else if (res.code == 200){
                setTimeout(function(){
                    parent.location.reload();
                }, 2000);
            } else {
                layer.msg("访问异常!", {icon: 2});
            }
        },
        error:function(e){
        	layer.close(index);
            layer.msg(e, {icon: 2});
        }
    });
    return false;
});

// 时间范围选择
laydate.render({
    elem: 'input[name="times"]',
    type: 'datetime',
    range: true,
    rangeLinked: true,
    trigger: 'click'
});

// 提示
$('.layui-input,.layui-textarea').click(function(){
    if($(this).attr('placeholder') != ''){
        var tps = $(this).attr('placeholder');
        if (tps.indexOf("*")>-1){
        } else {
            layer.tips(tps, $(this),{tips:1});
        }   
    }
});

//
$(document).ready(function(){
    var index;
    $('.tips').mouseenter(function(){
        var tps = $(this).attr('placeholder');
        if(tps != ''){
            index = layer.tips(tps, this, {
                tips: [1], // tips 方向和颜色
                time: 0 // 不自动关闭
            });
        }
    });

    $('.tips').mouseleave(function(){
        layer.close(index);
    });
});
///

///
});!function (win) {
///
"use strict";
var doc = document,
Admin = function(){
    this.v = '1.0'; //版本号
};
//默认加载
Admin.prototype.init = function () {
};

Admin.prototype.getRand = function(_id){
    var rand = Math.random().toString(36).substr(2)+Math.random().toString(36).substr(5);
    $('#'+_id).val(rand);
};

Admin.prototype.getPass = function(url){
    layer.prompt({title: '请输入新密码',area: ['200px', '150px']},function(value, index, elem){
        $.post(url, {pass:value}, function(res) {
            if(res.code == 1){
                layer.msg('修改成功',{icon: 1});
                setTimeout(function() {
                    location.reload();
                }, 1000);
            }else{
                layer.msg(res.msg,{icon: 2});
                layer.close(index);
            }
        },'json');
    });
}

//批量删除
Admin.prototype.batchDel = function(_url,_id) {
    var ids = [];
    if (isNaN(_id)) {
        var checkStatus = table.checkStatus(_id);
        checkStatus.data.forEach(function(n,i){
            ids.push(n.id);
        });
        var one = false;
    }else{
    	ids.push(_id);
    	var one = true;
    }

    if(ids.length == 0 ){
        layer.msg('请选择要删除的数据~!',{icon: 2,shift:6});
    }else{
        layer.confirm('确定要删除吗?', { title:'删除提示', btn: ['确定', '取消'],shade:0.001}, function(index) {
            $.post(_url, {'id':ids}, function(res) {
            	showMsg(res.msg, function(){
	        		if(res.code > -1){
	            		table.reload(_id);
	           	 	}
	            },{icon: res.code > -1 ? 1 : 2,shift:res.code ? 0 : 6});
            },'json');
        }, function(index) {
            layer.close(index);
        });
    }
};

Admin.prototype.del = function(_this,_url,_id) {
    layer.confirm('确定要删除吗?', { title:'删除提示', btn: ['确定', '取消'],shade:0.001}, function(index) {
        $.post(_url, { 'id':_id }, function(res) {
            console.log(res);
            showMsg(res.msg, function(){
        		if(res.code > -1){
            		location.reload();
           	 	}
            },{icon: res.code > -1 ? 1 : 2});
        },'json');
    }, function(index) {
        layer.close(index);
    });
};

Admin.prototype.trigger = function(_this,_url,_id) {
    layer.confirm('确定要触发吗?', { title:'触发提示', btn: ['确定', '取消'],shade:0.001}, function(index) {
        $.post(_url, { 'id':_id }, function(res) {
            showMsg(res.msg, function(){
                if(res.code > -1){
                    location.reload();
                }
            },{icon: res.code > -1 ? 1 : 2});
        },'json');
    }, function(index) {
        layer.close(index);
    });
};

//弹出层
Admin.prototype.open = function (title,url,w,h,full) {
    if (title == null || title == '') {
        var title = false;
    };
    if (w == null || w == '') {
        var w = ($(window).width()*0.9);
    };
    if (h == null || h == '') {
        var h = ($(window).height() - 50);
    };
    h = h-20;
    var open = layer.open({
        type: 2,
        area: [w+'px', h +'px'],
        fix: false, //不固定
        maxmin: true,
        shadeClose: false,
        maxmin: false,
        shade:0.2,
        title: title,
        content: url
    });
    if(full){
       layer.full(open);
    }
};
win.Admin = new Admin();
///
}(window);