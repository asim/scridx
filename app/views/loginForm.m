<div class="form-center">
<h1>Login</h1>
<form id="loginForm"  action="/login" method="POST">
 <fieldset>
  <div class="control-group control-top">
  <label class="control-label" for="inputUsername">Username</label>
   <div class="controls">
      <div class="holding">
       <input type="text" name="Username" id="inputUsername" placeholder="Username">
      </div>
   </div>
  </div>
  <div class="control-group">
   <label class="control-label" for="inputPassword">Password</label>
   <div class="controls">
      <div class="holding">
       <input type="password" name="Password" id="inputPassword" placeholder="Password">
      </div>
   </div>
  </div>
  <div class="control-group">
    <div class="controls">
     <label class="checkbox">
      <input type="checkbox" name="Remember" value="True"> Remember me
     </label>
     <button type="submit" class="btn btn-primary">Sign in</button>
   </div>
  </div>
 </fieldset>
 <input type="hidden" name="_csrf" value="{{csrf}}">
</form>
</div>
