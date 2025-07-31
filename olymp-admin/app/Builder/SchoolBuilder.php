<?php

namespace App\Builder;

use App\Models\School;

class SchoolBuilder
{
    public function build($item){
        $model = new School();
        $model->id = $item['id'];
        $model->name = $item['name'];
        $model->region = $item['region'];
        return $model;
    }
}
