import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import {Service} from "../../service"; // Assuming you have an AuthService to handle authentication

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html'
})
export class LoginComponent implements OnInit {
  username: string = '';
  password: string = '';

  constructor(private service: Service, private router: Router) { }

  ngOnInit(): void {
  }

  onLoginClick() {
    if (this.username && this.password) {
      this.service.login(this.username, this.password).subscribe(
        (response: any) => {
          const token = response.token;
          if (token) {
            localStorage.setItem('jwtToken', token);
            this.router.navigate(['/collections']);
          }
        }
      );
    }
  }
}
