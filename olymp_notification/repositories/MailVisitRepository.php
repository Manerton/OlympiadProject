<?php

namespace app\repositories;

use app\models\MailVisit;

class MailVisitRepository
{
    public function save(MailVisit $model)
    {
        if (!$model->save()) {
            throw new \RuntimeException('Saving error.');
        }
    }
}