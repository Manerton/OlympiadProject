<?php

namespace repositories;

use app\models\PhoneVisit;

class PhoneVisitRepository
{
    public function save(PhoneVisit $model)
    {
        if (!$model->save()) {
            throw new \RuntimeException('Saving error.');
        }
    }
}