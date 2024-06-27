import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params, Router } from '@angular/router';
import { Service } from 'src/app/service';

@Component({
  selector: 'app-new-item',
  templateUrl: './new-item.component.html',
})
export class NewItemComponent implements OnInit {
  constructor(private service: Service, private route: ActivatedRoute, private router: Router) { }

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

  createItem(title: string, rarity: string) {
    if (!title) {
      return;
    }

    const item = {
      title: title,
      rarity: Number(rarity)
    };

    this.service.createItem(item, this.collectionId, this.token).subscribe((newItem: any) => {
      this.router.navigate(['../'], { relativeTo: this.route });
    });
  }

}
