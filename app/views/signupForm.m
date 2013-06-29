<div class="form-center">
<h1>Signup</h1>
<form id="signupForm" action="/signup" method="POST">
 {{#form}}
 <fieldset>
  <div class="control-group control-top">
   <label class="control-label" for="inputName">Display Name</label>
   <div class="controls">
    <div class="holding"> 
      <input type="text" name="Name" id="inputName" placeholder="Name" value="{{Name}}">
    </div>
   </div>
  </div>
  <div class="control-group">
   <label class="control-label" for="inputUsername">Username</label>
   <div class="controls">
    <div class="holding"> 
      <input type="text" name="Username" id="inputUsername" placeholder="Username" value="{{Username}}">
    </div>
   </div>
  </div>

  <div class="control-group">
   <label class="control-label" for="inputEmail">Email</label>
   <div class="controls">
    <div class="holding"> 
     <input type="email" name="Email" id="inputEmail" placeholder="Email" value="{{Email}}">
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
  <button type="submit" class="btn btn-primary">Signup</button>
  </fieldset>
  {{/form}}
  <input type="hidden" name="_csrf" value="{{csrf}}">
</form>
</div>
