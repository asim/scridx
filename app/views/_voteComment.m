{{#Votable}}
<ul class="unstyled" id="ul_{{Id}}">
 <li><a href="{{VoteUrl}}/up" id="c_{{Id}}" {{#_user}}onclick="return vote(this)"{{/_user}}><i class="icon-chevron-up"></i></a></li>
 <li><a href="{{VoteUrl}}/down" id="c_{{Id}}" {{#_user}}onclick="return vote(this)"{{/_user}}><i class="icon-chevron-down"></i></a></li>
</ul>
{{/Votable}}
