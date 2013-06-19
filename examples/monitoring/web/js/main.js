require(["d3.v3.min", "ctl"], function() {
  require(["cubism.v1.min"], function() {
    var ctl = new Ctl();
    d3.json("/config", function(config) {
      ctl.config_cubism(config);
      ctl.poll("")
    });
  });
});
