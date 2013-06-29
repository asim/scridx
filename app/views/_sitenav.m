<div class="sitenav">
 <div class="sitenav-container">
 <div class="sitenav-body">
 <div id="sitenav-navbar">
  <div class="container-fluid">  
   <div class="navbar ">
    <div class="">
     <div class="container-fluid">
      <div class="logo-container">
       <div> 
        <a class="logo pull-left" href="/">
         <img alt="Scridx" src="/static/img/scridx_logo.png">
        </a>
       </div>
       <div>
        <a class="nav-btn" data-toggle="collapse" data-target=".nav-collapse">
         <img alt="Scridx" src="/static/img/scridx_icon.png">
        </a>
       </div>
      </div>

      <div class="nav-collapse collapse">
       <ul class="sitenav-ul">
        <li>
         <div><a href="/scripts">Scripts</a></div>
         <div><a href="/submit/script" class="btn btn-primary btn-mini">+ Add</a></div>
        </li>
        <li>
         <div><a href="/requests">Requests</a></div>
         <div><a href="/submit/request" class="btn btn-primary btn-mini">+ Ask</a></div>
        </li>
        <li>
         <div><a href="/feedback">Feedback</a></div>
         <div><a href="/submit/feedback" class="btn btn-primary btn-mini">+ Get</a></div>
        </li>
        <li>
         <div><a href="/news">News & Misc</a></div>
         <div><a href="/submit/news" class="btn btn-primary btn-mini">+ New</a></div>
        </li>
        </ul>
       </div> <!-- /nav-collapse -->
      </div> <!-- /container-fluid --> 
     </div> <!-- /navbar-inner -->
    </div> <!-- /navbar -->
  </div> <!-- /container-fluid -->
 </div> <!-- /sitenav-navbar -->
   {{#user}}
    <div class="sitenav-userinfo top-buffer">
    <div class="user-info">
     <ul class="unstyled">
      <li><img src="/static/img/av.png"/></li>
      <li><h3>{{Name}}</h3></li>
      <li class="user-logline"><p><small><em>{{Logline}}</em></small></p></li>
      {{#edit}}<li><a href="/settings">Edit</a>{{/edit}}
     </ul>
    </div>
    </div>
   {{/user}}
 </div> <!-- /sitenav-body -->

 <div class="sitenav-foot">
   <ul class="inline top-buffer">
    {{^_user}}
    <li><a href="/login">Login</a></li>
    <li><a href="/signup">Signup</a></li>
    {{/_user}}
    {{#_user}}
    <li><a href="/u/{{Username}}">Me</a></li>
    <li><a href="/logout">Logout</a></li>
    {{/_user}}
   </ul> 
  </div>
</div> <!-- /sitenav-container -->
</div> <!-- /sitenav -->
