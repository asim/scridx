<div class="form-center">
<h1>Request Script</h1>
<form id="requestForm" action="/submit/request" method="POST">
 <fieldset>
  <legend>Request Details</legend>
  {{#form}}

  <div class="control-group">
   <label class="control-label">Title</label>
   <div class="controls">
    <input class="span12" type=text name="Title" placeholder="Title" value="{{Title}}" autofocus>  
   </div>
  </div>

  <p>
   <a class="accordion-toggle" data-toggle="collapse" href="#collapseOne">
    [+] Optional (Writers, Draft Date, Logline, IMDB, Wikipedia)
   </a>
  </p>

  <div id="collapseOne" class="accordion-body collapse">
   <div class="control-group">
    <label class="control-label">Writer(s)</label>
    <div class="controls">
     <input class="span12" type=text name="Writers" placeholder="Writer 1, Writer 2, Writer..." value="{{Writers}}">
    </div>
   </div>
   <div class="control-group">
    <label class="control-label">IMDB Link</label>
    <div class="controls">
     <input class="span12" type=url name="Imdb" placeholder="IMBD Link" value="{{Imdb}}"> 
    </div>
   </div>
   <div class="control-group">
    <label class="control-label">Wikipedia Link</label>
    <div class="controls">
     <input class="span12" type=url name="Wiki" placeholder="Wiki Link" value="{{Wiki}}"> 
    </div>
   </div>
   <div class="control-group control-shared">
    <div class="controls control-left">
     <label class="control-label">Draft Date</label>
     <input class="span12" type=text name="Drafted" pattern="^(0[1-9]|1[012])\/(0[1-9]|[12][0-9]|3[01])\/(19|20)\d\d$" placeholder="MM/DD/YYYY" value="{{Drafted}}">
    </div>
    <div class="controls control-right">
     <label class="control-label">Draft Version</label>
     <input class="span12" type=text name="Version" placeholder="Version" value="{{Version}}">
    </div>
   </div>
   <div class="control-group">
    <label class="control-label">Logline</label>
    <div class="controls">
     <textarea name=Logline rows="5" class="span12" placeholder="Logline...">{{Logline}}</textarea> <br>
    </div>
   </div>
  </div>
   <button type=submit class="btn btn-primary">Submit</button>
  {{/form}}
 </fieldset>
 <input type="hidden" name="_csrf" value="{{csrf}}">
</form>
</div>
