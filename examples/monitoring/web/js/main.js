require(["d3.v3.min"], function() {
  require(["cubism.v1.min", "ctl"], function() {

    console.log(d3);
    console.log(cubism);
    var ctl = new Ctl();
    d3.json("/config", function(config) {
      ctl.config_cubism(config);
      ctl.poll("")
    });
  });
});
