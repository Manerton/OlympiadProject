<?php

namespace App\Services;

use App\Repositories\AppealRepository;

class AppealService
{
    private AppealRepository $appealRepository;
    public function __construct(
        AppealRepository $appealRepository
    )
    {
        $this->appealRepository = $appealRepository;
    }

    public function create($data)
    {
        $this->appealRepository->create($data['user_id'], $data['task_id'], $data['reason']);
    }
    public function changeStatus($id, $data){
        $model = $this->appealRepository->get($id);
        $this->appealRepository->changeStatus($model, $data['status']);
    }
}
