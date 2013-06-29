<div class="form-center">
<h1>Settings</h1>
<form id="settingsForm" action="/settings" method="POST">
 {{#user}}
 <fieldset>
  <div class="control-group control-top">
   <label class="control-label" for="inputName">Display Name</label>
   <div class="controls">
    <div class="holding"> 
      <input type="text" name="Name" id="inputName" placeholder="Display Name" value="{{Name}}">
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
   <label class="control-label" for="Logline">Logline (160 characters or less)</label>
   <div class="controls">
    <div class="holding"> 
     <textarea name=Logline rows="5" class="span12" placeholder="Logline...">{{Logline}}</textarea> <br>
    </div>
   </div>
  </div>
  <button type="submit" class="btn btn-primary">Save</button>
  </fieldset>
    
  {{/user}}
  <input type="hidden" name="_csrf" value="{{csrf}}">
</form>
</div>
