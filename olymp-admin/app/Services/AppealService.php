<?php

namespace App\Services;

use App\Repositories\AppealRepository;
use App\Repositories\TaskRepository;

class AppealService
{
    private AppealRepository $appealRepository;
    private TaskRepository $taskRepository;
    public function __construct(
        AppealRepository $appealRepository,
        TaskRepository $taskRepository
    )
    {
        $this->appealRepository = $appealRepository;
        $this->taskRepository = $taskRepository;
    }

    public function create($data)
    {
        $this->appealRepository->create($data['user_id'], $data['task_id'], $data['reason']);
    }
    public function changeStatus($id, $data){
        $model = $this->appealRepository->get($id);
        $this->appealRepository->changeStatus($model, $data['status']);
    }
    public function getByEventId($id){
        return $this->taskRepository->getByEventId($id)
            ->pluck('appeal')
            ->all();
    }
    public function getByUserId($id)
    {
        return $this->appealRepository->getByUserId($id);
    }
}
