  let Message = JSON.stringify({  "EN":[{ "PleaseSignIn":"Please sign in!", "CommentSend":"Your comment has been submitted","Success":"Successfully completed."}], 
                                  "TR":[{ "PleaseSignIn":"Lütfen giriş yapınız!", "CommentSend":"Yorumunuz başarıyla gönderildi.","Success":"Başarıyla tamamlandı."}] })
  const LanguageMessage = JSON.parse(Message); 

  const monthNames = ["January", "February", "March", "April", "May", "June",
    "July", "August", "September", "October", "November", "December"
  ];

  function rate(value,film,filmtype) {
    var formData = {
      value:value,
      film: film,
      filmtype: filmtype, 
    }
      $.ajax({
        type: "POST",
        url: "/backend/rate",
        data: formData,
      }).done(function(data){ 
          if(data.trim() == "add"){
            window.location.reload();
          }else if (data.trim() == "permission") {
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: LanguageMessage[LanguageSession][0].PleaseSignIn,
              showConfirmButton: false,
              timer: 1500
            });
          }
      });
      event.preventDefault();
    }

  function favorite(argument,filmtype) {
    if(argument == null){
      argument = 1;
    }
    var formData = {
      film: argument,
      filmtype: filmtype, 
    }
      $.ajax({
        type: "POST",
        url: "/backend/favorite",
        data: formData,
      }).done(function(data){ 
          if(data.trim() == "add"){
            if(document.getElementById("btnfavorite"+argument.toString()) !=undefined){
              document.getElementById("btnfavorite"+argument.toString()).innerHTML = '<i class="fas fa-trash"></i> Remove';
              document.getElementById("btnfavorite"+argument.toString()).classList.add("bg-danger");
              document.getElementById("btnfavorite"+argument.toString()).classList.remove("bg-warning");
            }

            if(document.getElementById("favoritelist"+argument.toString()) != undefined){
              //document.getElementById("watchedlist"+argument.toString()).innerHTML = 'Watchedlist <i class="fas fa-plus fa-sm  ml-1"></i>';
              //document.getElementById("watchedlist"+argument.toString()).classList.remove("btn-success");
              document.getElementById("favoritelist"+argument.toString()).classList.add("film_icon_active");
            }
          }else if (data.trim() == "remove") {
            if(document.getElementById("btnfavorite"+argument.toString()) !=undefined){
              document.getElementById("btnfavorite"+argument.toString()).innerHTML = '<i class="fas fa-clock"></i> Add';
              document.getElementById("btnfavorite"+argument.toString()).classList.remove("bg-danger");
              document.getElementById("btnfavorite"+argument.toString()).classList.add("bg-warning");
            }

            if(document.getElementById("favoritelist"+argument.toString()) != undefined){
              //document.getElementById("watchedlist"+argument.toString()).innerHTML = 'Watchedlist <i class="fas fa-plus fa-sm  ml-1"></i>';
              //document.getElementById("watchedlist"+argument.toString()).classList.remove("btn-success");
              document.getElementById("favoritelist"+argument.toString()).classList.remove("film_icon_active");
            }
          }else if (data.trim() == "permission") {
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: LanguageMessage[LanguageSession][0].PleaseSignIn,
              showConfirmButton: false,
              timer: 1500
            });
          }
      });
      event.preventDefault();
    }

 function follow(argument) {
    if(argument == null){
      argument = 1;
    }
    var formData = {
      follow_owner: argument,
    }
      $.ajax({
        type: "POST",
        url: "/backend/follow",
        data: formData,
      }).done(function(data){ 
          if(data.trim() == "add"){
            if(document.getElementById("follow"+argument.toString()) != undefined){
              document.getElementById("follow"+argument.toString()).innerHTML = '<i class="fas fa-check fa-sm  ml-1"></i> Unfollow';
              document.getElementById("follow"+argument.toString()).classList.remove("btn-danger");
              document.getElementById("follow"+argument.toString()).classList.add("btn-dark");
            }
          }else if (data.trim() == "remove") {
            if(document.getElementById("follow"+argument.toString()) != undefined){
              document.getElementById("follow"+argument.toString()).innerHTML = '<i class="fas fa-user-plus fa-sm  ml-1"></i> Follow';
              document.getElementById("follow"+argument.toString()).classList.remove("btn-dark");
              document.getElementById("follow"+argument.toString()).classList.add("btn-danger");
            }
          }else if (data.trim() == "permission") {
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: LanguageMessage[LanguageSession][0].PleaseSignIn,
              showConfirmButton: false,
              timer: 1500
            });
          }
      });
      //event.preventDefault();
    }
 
   function createlist() {
    Swal.fire({
      title: 'Create List',
      html: `<input type="text" id="list_title" class="swal2-input" placeholder="List Title">
    <select class="swal2-select" aria-label="Share Options" id="list_type">
      <option value="everyone" selected>Everyone</option>
      <option value="me">Only me</option>
    </select>`,
      confirmButtonText: 'Create List',
      focusConfirm: false,
      showCancelButton: true,
      preConfirm: () => {
        const list_title = Swal.getPopup().querySelector('#list_title').value
        const select = Swal.getPopup().querySelector('#list_type')
        const list_type = select.options[select.selectedIndex].value
        if (!list_title || !list_type) {
          Swal.showValidationMessage(`Please enter title and type!`)
        }
        return { list_title: list_title, list_type: list_type }
      }
    }).then((result) => {
          var formData = {
            list_title: result.value.list_title,
            list_type: result.value.list_type,
          }
          $.ajax({
              type: "POST",
              url: "/backend/createlist",
              data: formData,
            }).done(function(data){
              if(data.trim() == "add"){
                Swal.fire({
                  position: 'top-end',
                  icon: 'success',
                  title: LanguageMessage[LanguageSession][0].Success,
                  showConfirmButton: false,
                  timer: 1500
                });
              }else{
                Swal.fire({
                  position: 'top-end',
                  icon: 'warning',
                  title: LanguageMessage[LanguageSession][0].PleaseSignIn,
                  showConfirmButton: false,
                  timer: 1500
                });
              }
            })
    })
  }

  function watchlist(argument,filmtype) {
    if(argument == null){
      argument = 1;
    }
    var formData = {
      film: argument,
      filmtype: filmtype, 
    }
      $.ajax({
        type: "POST",
        url: "/backend/watchlist",
        data: formData,
      }).done(function(data){ 
          if(data.trim() == "add"){
            if(document.getElementById("coverwatchlist"+argument.toString()) != undefined){
              document.getElementById("coverwatchlist"+argument.toString()).classList.remove("text-white");
              document.getElementById("coverwatchlist"+argument.toString()).classList.add("text-primary");
            }

            if(document.getElementById("btnwatchlist"+argument.toString()) !=undefined){
              document.getElementById("btnwatchlist"+argument.toString()).innerHTML = '<i class="fas fa-trash"></i> Remove';
              document.getElementById("btnwatchlist"+argument.toString()).classList.add("bg-danger");
              document.getElementById("btnwatchlist"+argument.toString()).classList.remove("bg-warning");
            }

            if(document.getElementById("watchlist"+argument.toString()) != undefined){
              //document.getElementById("watchlist"+argument.toString()).innerHTML = 'Watchlist <i class="fas fa-check fa-sm  ml-1"></i>';
              //document.getElementById("watchlist"+argument.toString()).classList.remove("btn-primary");
              document.getElementById("watchlist"+argument.toString()).classList.add("film_icon_active");
            }

          }else if (data.trim() == "remove") {
            if(document.getElementById("coverwatchlist"+argument.toString()) != undefined){
               document.getElementById("coverwatchlist"+argument.toString()).classList.add("text-white");
               document.getElementById("coverwatchlist"+argument.toString()).classList.remove("text-primary");
            }

            if(document.getElementById("btnwatchlist"+argument.toString()) !=undefined){
              document.getElementById("btnwatchlist"+argument.toString()).innerHTML = '<i class="fas fa-clock"></i> Add';
              document.getElementById("btnwatchlist"+argument.toString()).classList.remove("bg-danger");
              document.getElementById("btnwatchlist"+argument.toString()).classList.add("bg-warning");
            }

            if(document.getElementById("watchlist"+argument.toString()) != undefined){
              //document.getElementById("watchlist"+argument.toString()).innerHTML = 'Watchlist <i class="fas fa-plus fa-sm  ml-1"></i>';
              //document.getElementById("watchlist"+argument.toString()).classList.remove("btn-success");
              document.getElementById("watchlist"+argument.toString()).classList.remove("film_icon_active");
            }
          }else if (data.trim() == "permission") {
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: LanguageMessage[LanguageSession][0].PleaseSignIn,
              showConfirmButton: false,
              timer: 1500
            });
            document.getElementById(argument).innerHTML = '<i class="fas fa-bookmark"></i> Add';
            document.getElementById(argument).classList.remove("bg-warning");
            document.getElementById(argument).classList.add("bg-danger");
          }
      });
      event.preventDefault();
    }

  function watchedlist(argument,filmtype,diary) {
    if(argument == null){
      argument = 1;
    }
    var formData = {
      film: argument,
      diary:diary,
      filmtype: filmtype, 
    }
    if(diary == null){
      $.ajax({
        type: "POST",
        url: "/backend/watchedlist",
        data: formData,
      }).done(function(data){ 
          if(data.trim() == "add"){
            if(document.getElementById("coverwatchedlist"+argument.toString()) != undefined){
              document.getElementById("coverwatchedlist"+argument.toString()).classList.remove("text-white");
              document.getElementById("coverwatchedlist"+argument.toString()).classList.add("text-primary");
            }

            if(document.getElementById("btnwatchedlist"+argument.toString()) !=undefined){
              document.getElementById("btnwatchedlist"+argument.toString()).innerHTML = '<i class="fas fa-trash"></i> Remove';
              document.getElementById("btnwatchedlist"+argument.toString()).classList.add("bg-danger");
              document.getElementById("btnwatchedlist"+argument.toString()).classList.remove("bg-warning");
            }

            if(document.getElementById("watchedlist"+argument.toString()) != undefined){
              //document.getElementById("watchedlist"+argument.toString()).innerHTML = 'Watchedlist <i class="fas fa-check fa-sm  ml-1"></i>';
              //document.getElementById("watchedlist"+argument.toString()).classList.remove("btn-danger");
              document.getElementById("watchedlist"+argument.toString()).classList.add("film_icon_active");

              document.getElementById("diary"+argument.toString()).classList.add("d-sm-inline-block");
              document.getElementById("diary"+argument.toString()).classList.remove("d-none");
            }
          }else if (data.trim() == "remove") {
            if(document.getElementById("coverwatchedlist"+argument.toString()) != undefined){
              document.getElementById("coverwatchedlist"+argument.toString()).classList.add("text-white");
              document.getElementById("coverwatchedlist"+argument.toString()).classList.remove("text-primary");
            }

            if(document.getElementById("btnwatchedlist"+argument.toString()) !=undefined){
              document.getElementById("btnwatchedlist"+argument.toString()).innerHTML = '<i class="fas fa-bookmark"></i> Add';
              document.getElementById("btnwatchedlist"+argument.toString()).classList.remove("bg-danger");
              document.getElementById("btnwatchedlist"+argument.toString()).classList.add("bg-warning");
            }

            if(document.getElementById("watchedlist"+argument.toString()) != undefined){
              //document.getElementById("watchedlist"+argument.toString()).innerHTML = 'Watchedlist <i class="fas fa-plus fa-sm  ml-1"></i>';
              //document.getElementById("watchedlist"+argument.toString()).classList.remove("btn-success");
              document.getElementById("watchedlist"+argument.toString()).classList.remove("film_icon_active");
              
              document.getElementById("diary"+argument.toString()).classList.remove("d-sm-inline-block");
              document.getElementById("diary"+argument.toString()).classList.add("d-none");
              document.getElementById("diary"+argument.toString()).classList.remove("d-block");
            }
          }else if (data.trim() == "permission") {
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: LanguageMessage[LanguageSession][0].PleaseSignIn,
              showConfirmButton: false,
              timer: 1500
            });
          }
      });
      event.preventDefault();
      }else{
        $.ajax({
          type: "POST",
          url: "/backend/diary",
          data: formData,
        }).done(function(data){ 
            if(data.trim() == "add"){
              if(document.getElementById("btndiary"+argument.toString()) !=undefined){
                document.getElementById("btndiary"+argument.toString()).innerHTML = '<i class="fas fa-trash"></i> Remove';
                document.getElementById("btndiary"+argument.toString()).classList.add("bg-danger");
                document.getElementById("btndiary"+argument.toString()).classList.remove("bg-warning");
              }

              if(document.getElementById("diary"+argument.toString()) != undefined){
                document.getElementById("diary"+argument.toString()).classList.add("film_icon_active");
              }
            }else if (data.trim() == "remove") {
              if(document.getElementById("btndiary"+argument.toString()) !=undefined){
                document.getElementById("btndiary"+argument.toString()).innerHTML = '<i class="fa-regular fa-calendar-days"></i> Add';
                document.getElementById("btndiary"+argument.toString()).classList.remove("bg-danger");
                document.getElementById("btndiary"+argument.toString()).classList.add("bg-warning");
              }

              if(document.getElementById("diary"+argument.toString()) != undefined){
                document.getElementById("diary"+argument.toString()).classList.remove("film_icon_active");
              }
            }else if (data.trim() == "permission") {
              Swal.fire({
                position: 'top-end',
                icon: 'warning',
                title: LanguageMessage[LanguageSession][0].PleaseSignIn,
                showConfirmButton: false,
                timer: 1500
              });
            }
        });
        event.preventDefault();
      }
    }

    function Language(argument) {
        let date = new Date();
        let exp = date.setTime(date.getTime() + (30 * 24 * 60 * 60 * 1000));
        document.cookie = "Language=" + argument + "; " + exp + "; path=/";

        if(argument=="Turkish"){
          const http = new XMLHttpRequest();
          const url='/backend/language-'+"tr";
          http.open("GET", url);
          http.send();
          http.onreadystatechange = (e) => {console.log('done')}
        }else{
          const http = new XMLHttpRequest();
          const url='/backend/language-'+"en";
          http.open("GET", url);
          http.send();
          http.onreadystatechange = (e) => {console.log('done')}
        }

        if(argument=="English"){
          var ar="English"
        }else if(argument=="Turkish"){
          var ar="Turkish"
        }
              Swal.fire({
              position: 'top-end',
              icon: 'success',
              title: 'You selected '+ar,
              showConfirmButton: false,
              timer: 1500
            });
              
            setTimeout(function(){window.location.reload();}, 1500);
    }


$( document ).ready(function() {
  const d = new Date();
  let time = d.getTime();

  $(".carousel-item").first().addClass("active"); 

  var current = window.location.href.replace("&in=movies","").replace("&in=people","").replace("&in=tvshow","");
  if(window.location.href.includes("/search")){
    document.getElementById("inmovies").setAttribute('href', current+"&in=movies"); 
    document.getElementById("inpeople").setAttribute('href', current+"&in=people"); 
    document.getElementById("intvshow").setAttribute('href', current+"&in=tvshow"); 
    
    if(window.location.href.includes("&in=people")){
      document.getElementById("inpeople").classList.add("active");
    }else if(window.location.href.includes("&in=tvshow")){
      document.getElementById("intvshow").classList.add("active");
    }else{
       document.getElementById("inmovies").classList.add("active");
    }
  }


  if(document.getElementById("content") != undefined){
      var input = document.getElementById("content");
    input.addEventListener("keypress", function(event) {
      if (event.key === "Enter") {
        event.preventDefault();
        if(document.getElementById("sendbutton") != undefined){document.getElementById("sendbutton").click();}
      }
    });
  }




  $("#SignIn").submit(function (event) {
    var formData = {
      kullaniciadi: $("#kullaniciadi").val(),
      sifre: $("#sifre").val(),
    };

    $.ajax({
      type: "POST",
      url: "/signin",
      data: formData,
    }).done(function(data){ 
        if(data.trim() === "ok"){
          setTimeout(function(){window.location="/";}, 500);
          Swal.fire({
            position: 'top-end',
            icon: 'success',
            title: 'Hesabınıza giriş yaptınız..',
            showConfirmButton: false,
            timer: 1000
          })

        }else if(data.trim() === "error"){
          Swal.fire({
            position: 'top-end',
            icon: 'error',
            title: 'Yanlış kullanıcı adı veya şifre!',
            showConfirmButton: false,
            timer: 1000
          })
        }
    });
    event.preventDefault();
  });


  $("#SignUp").submit(function (event) {
      var formData = {
        nick: $("#nick").val(),
        name: $("#name").val(), 
        email: $("#email").val(),
        password1: $("#password1").val(),
        password2: $("#password2").val(),
      }

      $.ajax({
        type: "POST",
        url: "/signup",
        data: formData,
      }).done(function(data){ 
          if(data.trim() === "ok"){

            setTimeout(function(){window.location="/signin";}, 2500);
            Swal.fire({
              position: 'top-end',
              icon: 'success',
              title: 'Mailinize gönderilen linke tıklayınız.',
              showConfirmButton: false,
              timer: 2500
            })
          }else if(data.trim() === "insecure"){
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: 'Insecure password, try including more special characters, using uppercase letters, using numbers or using a longer password',
              showConfirmButton: false,
              timer: 1500
            })
          }else if(data.trim() === "password"){
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: 'İlk yazdığınız şifre ile ikincisi uyuşmuyor!',
              showConfirmButton: false,
              timer: 1500
            })
          }else if(data.trim() === "nick"){
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: 'Bu kullanıcı adını (nick) başka bir kullanıcı kullanıyor!',
              showConfirmButton: false,
              timer: 1500
            })

          }else if(data.trim() === "email"){
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: 'Bu e-postayı başka bir kullanıcı kullanıyor!',
              showConfirmButton: false,
              timer: 1500
            })
          }else if(data.trim() === "error"){
            Swal.fire({
              position: 'top-end',
              icon: 'error',
              title: 'Bir şeyleri eksik bırakmış olmalısınız!',
              showConfirmButton: false,
              timer: 1500
            })
          }
      });
      event.preventDefault();
  });

  $("#ChangeSettings").submit(function (event) {
    var formData = {
      name: $("#name").val(),
      telephone: $("#telephone").val(),
      email: $("#email").val(),
      bio: $("#bio").val(),
    };

    $.ajax({
      type: "POST",
      url: "/backend/settings",
      data: formData,
    }).done(function(data){ 
        if(data.trim() === "ok"){
          setTimeout(function(){window.location.reload()}, 1000);
          Swal.fire({
            position: 'top-end',
            icon: 'success',
            title: 'Settings changed!',
            showConfirmButton: false,
            timer: 1500
          })
        }else if(data.trim() === "mail"){
          Swal.fire({
            position: 'top-end',
            icon: 'warning',
            title: 'mail alert',
            showConfirmButton: false,
            timer: 1500
          })
        }
    });
    event.preventDefault();
  });

  $("#ChangePassword").submit(function (event) {
    var formData = {
      oldpassword: $("#oldpassword").val(),
      newpassword1: $("#newpassword1").val(),
      newpassword2: $("#newpassword2").val(),
    };

    $.ajax({
      type: "POST",
      url: "/backend/changepassword",
      data: formData,
    }).done(function(data){ 
        if(data.trim() === "ok"){
          setTimeout(function(){window.location.reload()}, 1000);
          Swal.fire({
            position: 'top-end',
            icon: 'success',
            title: 'Password changed!',
            showConfirmButton: false,
            timer: 1500
          })
        }else if(data.trim() === "empty"){
          Swal.fire({
            position: 'top-end',
            icon: 'warning',
            title: 'empty',
            showConfirmButton: false,
            timer: 1500
          })
        }else if(data.trim() === "newpassword"){
          Swal.fire({
            position: 'top-end',
            icon: 'warning',
            title: 'New password1 and new password2 are not the same!',
            showConfirmButton: false,
            timer: 1500
          })
        }else if(data.trim() === "oldpassword"){
          Swal.fire({
            position: 'top-end',
            icon: 'warning',
            title: 'Old password is wrong!',
            showConfirmButton: false,
            timer: 1500
          })
        }else if(data.trim() === "insecure"){
          Swal.fire({
            position: 'top-end',
            icon: 'warning',
            title: 'Insecure new password, try including more special characters, using uppercase letters, using numbers or using a longer password',
            showConfirmButton: false,
            timer: 1500
           })
        }
    });
    event.preventDefault();
  });

 $("#messagesend").submit(function (event) {
      var formData = {
        target: $("#target").val(),
        content: $("#content").val(), 
      }
      $.ajax({
        type: "POST",
        url: "/backend/messagesend",
        data: formData,
      }).done(function(data){ 
          if(data.trim()==="ok"){
            document.getElementById("content").value = "";

            $('#ChatBox').stop ().animate ({
              scrollTop: $('#ChatBox')[0].scrollHeight
            });

          }else if(data.trim() === "content"){}
      });
      event.preventDefault();
  });

  $("#FilmComment").submit(function (event) {
      var formData = {
        film: $("#film").val(),
        content: $("#ReviewContent").val(),
        filmtype: $("#filmtype").val(),
      }

      $.ajax({
        type: "POST",
        url: "/backend/commentfilm",
        data: formData,
      }).done(function(data){ 
          if(data.trim()==="ok"){
            setTimeout(function(){window.location.reload()}, 1500);
            Swal.fire({
              position: 'top-end',
              icon: 'success',
              title: LanguageMessage[LanguageSession][0].CommentSend,
              showConfirmButton: false,
              timer: 1500
            });

          }else if(data.trim() === "film"){
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: 'Film verisi çekilemedi! Sayfayı yenileyiniz.',
              showConfirmButton: false,
              timer: 1500
            })
          }else if(data.trim() === "content"){
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: 'Yorumunuz boş olamaz!',
              showConfirmButton: false,
              timer: 1500
            });

          }
      });
      event.preventDefault();
  });


  $("#CommentBlog").submit(function (event) {
      var formData = {
        blog: $("#blog").val(),
        content: $("#content").val(), 
        rate: $("#rate").val(),
      }

      $.ajax({
        type: "POST",
        url: "/backend/commentblog",
        data: formData,
      }).done(function(data){ 
          if(data.trim().includes("ok")){
            setTimeout(function(){window.location.reload();}, 1500);
            Swal.fire({
              position: 'top-end',
              icon: 'success',
              title: 'Yorumunuz gönderildi!',
              showConfirmButton: false,
              timer: 1500
            })

          }else if(data.trim() === "blog"){
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: 'Blog verisi çekilemedi! Sayfayı yenileyiniz.',
              showConfirmButton: false,
              timer: 1500
            })
          }else if(data.trim() === "content"){
            Swal.fire({
              position: 'top-end',
              icon: 'warning',
              title: 'Yorumunuz boş olamaz!',
              showConfirmButton: false,
              timer: 1500
            })

          }
      });
      event.preventDefault();
  });

});

function CopyUrl(){
  var inputc = document.body.appendChild(document.createElement("input"));
  inputc.value = window.location.href;
  inputc.select();
  document.execCommand('copy');
  inputc.parentNode.removeChild(inputc);
            Swal.fire({
              position: 'top-end',
              icon: 'success',
              title: 'Copied to clipboard!',
              showConfirmButton: false,
              timer: 1000
            })
}

function escape_html(str) {
 if ((str===null) || (str===''))
       return false;
 else
   str = str.toString();
  var map = {
    '&': '&amp;',
  '<': '&lt;',
  '>': '&gt;',
  '"': '&quot;',
  "'": '&#039;'
  };

  return str.replace(/[&<>"']/g, function(m) { return map[m]; });
}


