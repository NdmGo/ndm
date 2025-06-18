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

function toSize(a) {
    var d = [" B", " KB", " MB", " GB", " TB", " PB"];
    var e = 1024;
    for(var b = 0; b < d.length; b++) {
        if(a < e) {
            return(b == 0 ? a : a.toFixed(2)) + d[b]
        }
        a /= e;
    }
}

function toTrim(x) {
    return x.replace(/^\s+|\s+$/gm,'');
}

function inArray(f, arr){
    for (var i = 0; i < arr.length; i++) {
        if (f == arr[i]) {
            return true;
        }
    }
    return false;
}


function isoTimeFormat(isoDateStr){
    // const isoDateStr = "2025-05-30T00:51:42.721+08:00";

    // 1. 解析 ISO 时间字符串
    const date = new Date(isoDateStr);

    // 2. 提取年月日时分秒
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, "0"); // 0-11 → 补0
    const day = String(date.getDate()).padStart(2, "0");
    const hours = String(date.getHours()).padStart(2, "0");
    const minutes = String(date.getMinutes()).padStart(2, "0");
    const seconds = String(date.getSeconds()).padStart(2, "0");

    // 3. 组合成目标格式
    const formattedDate = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
    // console.log(formattedDate); // 输出: "2025-05-30 00:51:42"
    return formattedDate
}

function getFileTypeFormat(b){
    switch(b) {
        case "html":
            var j = {
                name: "htmlmixed",
                scriptTypes: [{
                    matches: /\/x-handlebars-template|\/x-mustache/i,
                    mode: null
                }, {
                    matches: /(text|application)\/(x-)?vb(a|script)/i,
                    mode: "vbscript"
                }]
            };
            d = j;
            break;
        case "htm":
            var j = {
                name: "htmlmixed",
                scriptTypes: [{
                    matches: /\/x-handlebars-template|\/x-mustache/i,
                    mode: null
                }, {
                    matches: /(text|application)\/(x-)?vb(a|script)/i,
                    mode: "vbscript"
                }]
            };
            d = j;
            break;
        case "js":
            d = "text/javascript";
            break;
        case "json":
            d = "application/ld+json";
            break;
        case "css":
            d = "text/css";
            break;
        case "php":
            d = "application/x-httpd-php";
            break;
        case "py":
            d = "text/x-cython";
        case "sh":
            d = "shell";
        case "tpl":
            d = "application/x-httpd-php";
            break;
        case "xml":
            d = "application/xml";
            break;
        case "sql":
            d = "text/x-sql";
            break;
        case "conf":
            d = "text/x-nginx-conf";
            break;
        default:
            var j = {
                name: "htmlmixed",
                scriptTypes: [
                    {matches: /\/x-handlebars-template|\/x-mustache/i,mode: null}, 
                    {matches: /(text|application)\/(x-)?vb(a|script)/i,mode: "vbscript"}
                ]
            };
            d = j;
    }
    return d;
}

function getExtName(fileName){
    var extArr = fileName.split(".");
    var extLastName = extArr[extArr.length - 1];
    return extLastName;
}

function getFileExtension(filename) {
    return filename.split('.').pop();
}


function isImageFile(ext){
    if (inArray(ext,['png','jpeg','jpg','gif','webp','bmp','ico'])){
        return true;
    }
    return false;
}

function isVideoFile(ext){
    if (inArray(ext,['mp4','avi','mov','mkv','wmv','flv', 'webm','m4v','mpeg','3gp','rmvb','rm','m2ts'])){
        return true;
    }
    return false;
}

function isM3u8File(ext){
    if (inArray(ext,['m3u8'])){
        return true;
    }
    return false;
}

function isCodeFile(ext){
    if (inArray(ext,['html','htm','php','txt','md','js','css','scss','json','c','h','pl','py','java','log','conf','sh','json','ini', 'yaml'])){
        return true;
    }
    return false;
}


function filterAddition(name , val){
    if (name == 'show_hidden') {
        if (val == 'on'){
            return true;
        } else {
            return false;
        }
    } else if (name == 'enable_sign') {
        if (val == 'on'){
            return true;
        } else {
            return false;
        }
    }
    return false;
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
        xhr.setRequestHeader('Authorization', localStorage.getItem('token'));
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


Admin.prototype.Copy = function (text) {
    var textarea = document.createElement('textarea');
    textarea.value = text;
    document.body.appendChild(textarea);
    textarea.select();
    try {
        var successful = document.execCommand('copy');
        if(successful) {
            layer.msg('复制成功');
        } else {
            layer.msg('复制失败');
        }
    } catch (err) {
        layer.msg('复制错误: ' + err);
    }
    document.body.removeChild(textarea);
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

//Item渲染
Admin.prototype.renderFormItem = function (obj) {
    var html = '';
    for (var oi = 0; oi < obj.length; oi++) {
        var data = obj[oi];
        if (data['tag'] == 'input'){
            html += '<div class="layui-form-item">\
                <label class="layui-form-label">'+data['label']+'</label>\
                <div class="layui-input-block">\
                    <input type="'+data['type']+'" name="'+data['key']+'" value="'+data['default']+'" autocomplete="off" class="layui-input" placeholder="'+data['placeholder']+'">\
                </div>\
            </div>';
        } else if (data['tag'] == 'textarea'){
            html +=  '<div class="layui-form-item layui-form-text">\
                <label class="layui-form-label">'+data['label']+'</label>\
                <div class="layui-input-block">\
                    <textarea class="layui-textarea" name="'+data['key']+'">'+data['default']+'</textarea>\
                </div>\
            </div>';
        } else if ( data['tag'] == 'select'){
            var con = '';
            var option = data['option'];
            for (var ii = 0; ii < option.length; ii++) {
                con += '<option value="'+option[ii]['value']+'">'+option[ii]['name']+'</option>';
            }

            html += '<div class="layui-form-item">\
                    <label class="layui-form-label">'+data['label']+'</label>\
                    <div class="layui-input-block">\
                        <select name="'+data['key']+'">'+con+'</select>\
                    </div>\
                </div>';
        } else if ( data['tag'] == 'input_switch'){
            html +=  '<div class="layui-form-item">\
                    <label class="layui-form-label">启用签名</label>\
                    <div class="layui-input-block">\
                        <input type="checkbox" name="'+data['key']+'" lay-skin="switch" title="开关" '+((data['checked']) ? 'checked':'')+'>\
                    </div>\
                </div>';
        } else if ( data['tag'] == 'line'){
            html +=  '<fieldset class="layui-elem-field layui-field-title" style="text-align: center;">\
                    <legend style="font-size: 12px;">'+data['title']+'</legend>\
                </fieldset>';
        }
    }

    return html;
};

//Item渲染多选
Admin.prototype.renderFormItemMulti = function (objs, selected) {
    console.log(objs,selected);
    console.log(objs[selected]);
    return this.renderFormItem(objs[selected]);
};




win.Admin = new Admin();
///
}(window);