<!DOCTYPE html>
<html>
  <head>
    <title>Scridx</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <!-- Bootstrap -->
    <link href="/static/css/theme.min.css" rel="stylesheet" media="screen">
    <link href="/static/css/responsive.min.css" rel="stylesheet">

    <!-- HTML5 shim, for IE6-8 support of HTML5 elements -->
    <!--[if lt IE 9]>
      <script src="http://html5shim.googlecode.com/svn/trunk/html5.js"></script>
    <![endif]-->
  </head>
  <body>
  <div class="container-fluid">
   <div class="row-fluid container-body">
    <!-- nav -->
    {{> _sitenav.m}}
   <!-- /nav -->
   <!-- main -->
   <div class="main">
     <div id="body">
      <div class="container-fluid">
       <div id="wrap">
        <div id="main">
         {{> _flash.m}}
         {{{content}}}
        </div>
       </div>
      </div> <!-- /container -->
     </div>

     <!-- foot -->
     <div id="footer">
      <div class="container-fluid">
       {{> _pager.m}}
       <hr>
       <div class="pull-right"><a href="#">Back to top</a></div>
       <div class="pull-left">&copy; 2013 Scridx</div>
      </div>
     </div> 
     <!-- /foot -->
    </div> <!-- /span7 -->
    <!-- /main -->

   </div> <!-- /row-fluid -->
  </div> <!-- /container-fluid -->
  </body>
  <script src="//ajax.googleapis.com/ajax/libs/jquery/2.0.0/jquery.min.js"></script>
  <script src="/static/js/jquery.validate.min.js"></script>
  <script src="/static/js/bootstrap.min.js"></script>
  <script src="/static/js/j.js"></script>
  <script type="text/javascript">
    Scridx.init();
    $('.text:empty').hide();
    $('.things:empty').hide();
  </script>
 </html>

</html>
