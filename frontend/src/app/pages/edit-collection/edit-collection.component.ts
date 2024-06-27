import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Service } from 'src/app/service';

@Component({
  selector: 'app-edit-list',
  templateUrl: './edit-collection.component.html',
})
export class EditCollectionComponent implements OnInit{
  selectedListId: any;
  constructor(private route: ActivatedRoute, private service: Service, private router: Router) { }

  collectionId: any;
  token: any;


  ngOnInit() {
    this.token = localStorage.getItem('jwtToken') as string;

    if (this.token === null) {
      this.router.navigate(['/login']);
    }

    this.route.params.subscribe(
      (params: Params) => {
        this.collectionId = params['collectionId'];
      }
    )
  }

  updateCollection(title: string) {
    this.service.updateCollection(this.collectionId, title, this.token).subscribe((data) => {
      this.router.navigate(['/collections', this.collectionId]);
    })
  }
}
