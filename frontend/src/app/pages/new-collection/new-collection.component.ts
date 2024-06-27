import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Service } from 'src/app/service';

@Component({
  selector: 'app-new-list',
  templateUrl: './new-collection.component.html',
})
export class NewCollectionComponent implements OnInit{
    token: any;
    constructor(private service: Service, private router: Router){

  }

  ngOnInit() {
    this.token = localStorage.getItem('jwtToken') as string;

    if (this.token === null) {
      this.router.navigate(['/login']);
    }
  }
  createCollection(title: string){
    this.service.createCollection(title, this.token).subscribe((response: any)=>{
      this.router.navigate(["/collections", response.id]);
    });
  }
}
