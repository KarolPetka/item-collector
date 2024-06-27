import {Component, OnInit} from '@angular/core';
import {Router} from '@angular/router';
import {Service} from "../../service"; // Assuming you have an AuthService to handle authentication

@Component({
    selector: 'app-login',
    templateUrl: './signup.component.html',
})
export class SignupComponent implements OnInit {
    username: string = '';
    password: string = '';

    constructor(private service: Service, private router: Router) {
    }

    ngOnInit(): void {
    }

    onLoginClick() {
        if (this.username && this.password) {
            this.service.signup(this.username, this.password).subscribe(
                (response: any) => {
                  if (response){
                    this.router.navigate(['/login']);
                  }
                }
            )
        }
    }
}
