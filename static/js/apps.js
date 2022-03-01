layui.use(['element', 'layedit', 'form', 'layer'], function(){
    layui.form.on("submit(install)", function(data){
        layui.$.post("/install", data.field, function(){
            console.log(data.field)
        })
    })
    
})