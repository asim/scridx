function vote(n) {
  n.style.visibility = 'hidden';
  if (n.id != "") {
    var v = n.id.split(/_/);
    document.getElementById("ul_" + v[1]).style.visibility = 'hidden'; 
  }
  var ping = new Image();
  ping.src = n.href;
  return false;
};

var Scridx = {
  init: function () {
    var self = this;
    $.validator.addMethod("alphanumeric", function(value, element) {
      return this.optional(element) || /^[a-zA-Z0-9_]+$/.test(value);
    }); 

    $.validator.addMethod("alphanumericExtended", function(value, element) {
      return this.optional(element) || /^[A-Za-z0-9\u00C0-\u017F\s@&.,_-]+$/.test(value);
    }); 

    self.validateCommentForm();
    self.validateItemForm();
    self.validateSettingsForm();
    self.validateScriptForm();
    self.validateRequestForm();
    self.validateApproveRequestForm();
    self.validateFulfillRequestForm();
    self.validateSignupForm();
    self.validateLoginForm();
    self.vote;
    return self;
  },


  validateRules: function() {
    return {
     rules: {
      Title: {
        required: true,
        rangelength: [1, 255]
      },
      Source: {
        required: true,
        rangelength: [12, 255],
        url: true
      },
      DraftDate: {
        rangelength: [10, 10],
        required: false,
        date: true
      },
      Logline: {
        maxlength: 512,
        required: false
      },
      Imdb: {
        rangelength: [12, 255],
        required: false,
        url: true
      },
      Wiki: {
        rangelength: [12, 255],
        required: false,
        url: true
      },
      Writers: {
        rangelength: [2, 512],
        required: false
      },
      Version: {
        rangelength: [1, 20],
        required: false
      }
     },
     messages: {
       Title: "A script needs a title.",
       Source: "Provide the source pdf url.",
       DraftDate: "Please enter a valid date, MM/DD/YYYY.",
       Writers: "Provide a comma-seperated list e.g Aaron Sorkin, Josh Whedon, ...",
       Version: "Version should be less than 20 characters."
     },
     highlight: function(element) {
       $(element).closest('.control-group').removeClass('success').addClass('error');
     },
     success: function(element) {
       element
       .addClass('valid')
       .closest('.control-group').removeClass('error').addClass('success');
     }
    }
  },


  validateUserBase: function() {
    return {
      rules: {
        Username: {
          required: true,
          rangelength: [1, 15],
          alphanumeric: true
        },
        Password: {
          required: true
        }
      },
      messages: {
       Username: "Um, that doesn't look right.",
       Password: "Need a valid password.",
      },
      highlight: function(element) {
       $(element).closest('.control-group').removeClass('success').addClass('error');
      },
      success: function(element) {
       element
       .addClass('valid')
       .closest('.control-group').removeClass('error').addClass('success');
      },
      errorElement: "p",
      wrapper: "div",
        errorPlacement: function(error, element) {
            offset = element.offset();
            error.insertAfter(element)
            error.addClass('sidetip');  // add a class to the wrapper
            error.css('left', offset.left + element.outerWidth());
            error.css('top', offset.top);
        }

    }
  },

  validateCommentForm: function() {
    var base = Scridx.validateRules();

    var rules = {
     rules: {
      Text: {
        maxlength: 10000,
        required: true
      }
     },
     messages: {
       Text: {
         maxlength: "Max length 10000 characters.",
         required: "Comment cannot be blank."
       }
     }
    };

    var validateRules = $.extend(true, {}, base, rules)
    $('#commentForm').validate(validateRules);
  },


  validateLoginForm: function() {
    $('#loginForm').validate(Scridx.validateUserBase());
  },

  validateSignupForm: function() {
    var base = Scridx.validateUserBase();

    var rules = {
     rules: {
       Name: {
         required: true,
         rangelength: [1, 20],
       },
       Email: {
        required: true,
        email: true,
        maxlength: 254
      }
     },
     messages: {
       Name: "Got a name?",
       Username: "Alphanumeric characters only!",
       Password: "Use 8 characters or more.",
       Email: "Is that an email address?",
     }
    };

    var validateRules = $.extend(true, {}, base, rules)
    $('#signupForm').validate(validateRules);
  },

  validateFulfillRequestForm: function() {
    $('#fulfillRequestForm').validate(Scridx.validateRules());
  },
   
  validateApproveRequestForm: function() {
    $('#approveRequestForm').validate(Scridx.validateRules());
  },

  validateItemForm: function() {
    var base = Scridx.validateRules();

    var rules = {
     rules: {
       Title: {
         required: true,
         rangelength: [1, 255]
       },
      Source: {
        required: false,
        rangelength: [10, 255],
        url: true
      },
      Text: {
        maxlength: 10000,
        required: false
      }
     },
     messages: {
       Title: "Title is required and must be less than 255 characters.",
       Url: "Must be a valid url",
       Text: "Needs to be less than 10,000 words",
     }
    };

    var validateRules = $.extend(true, {}, base, rules)
    $('#itemForm').validate(validateRules);
  },


  validateRequestForm: function() {
    var validateRules = Scridx.validateRules();
    validateRules.rules.Source.required = false;

    $('#requestForm').validate(validateRules);
  },

  validateScriptForm: function() {
    $('#scriptForm').validate(Scridx.validateRules());
  },

  validateSettingsForm: function() {
    var base = Scridx.validateRules();

    var rules = {
     rules: {
       Name: {
         required: true,
         rangelength: [1, 20],
         alphanumericExtended: true
       },
       Email: {
        required: true,
        email: true,
        maxlength: 254
      },
      Logline: {
        required: false,
        maxlength: 160
      }
     },
     messages: {
       Name: "Is that right?",
       Email: "Is that an email address?",
       Logline: "Needs to be less than 160 characters",
     }
    };

    var validateRules = $.extend(true, {}, base, rules)
    $('#settingsForm').validate(validateRules);
  },
};
