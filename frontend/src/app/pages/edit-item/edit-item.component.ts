import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Service } from 'src/app/service';

@Component({
  selector: 'app-edit-item',
  templateUrl: './edit-item.component.html',
})
export class EditItemComponent implements OnInit{
  constructor(private route: ActivatedRoute, private service: Service, private router: Router) { }

  itemId: any;
  collectionId: any;
  token: any;


  ngOnInit() {
    this.token = localStorage.getItem('jwtToken') as string;

    if (this.token === null) {
      this.router.navigate(['/login']);
    }

    this.route.params.subscribe(
      (params: Params) => {
        this.itemId = params['itemId'];
        this.collectionId = params['collectionId'];
      }
    )
  }

  updateItem(title: string, rarity: string) {
    if (!title) {
      return;
    }

    const item = {
      title: title,
      rarity: Number(rarity)
    };

    this.service.updateItem(item, this.collectionId, this.itemId, this.token).subscribe(() => {
      this.router.navigate(['/collections', this.collectionId]);
    })
  }
}
