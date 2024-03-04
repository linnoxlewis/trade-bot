 #!/bin/bash

 # Переменная с путем к вашему проекту
 PROJECT_PATH=../internal

 # Путь к папке с репозиториями
 REPOSITORY_PATH="${PROJECT_PATH}/service"


# Проверяем, установлен ли mockgen
if ! command -v mockgen &> /dev/null; then
    echo "mockgen не установлен. Установите его, выполнив: go get github.com/golang/mock/mockgen"
    exit 1
fi

# Создаем моки для всех интерфейсов в папках repository
for folder in "${REPOSITORY_PATH}"/*; do
    if [ -d "${folder}" ]; then
        PACKAGE_NAME=$(basename "${folder}")

        # Создаем папку mocks, если ее нет
        MOCKS_DIR="${REPOSITORY_PATH}/${PACKAGE_NAME}/mocks"
        mkdir -p "${MOCKS_DIR}"

        # Генерация мока
        mockgen -destination="${MOCKS_DIR}/mock_${PACKAGE_NAME}.go" -package="mocks" "${PROJECT_PATH}/internal/service/${PACKAGE_NAME}" YourInterfaceToMock
        echo "Создан мок для ${PACKAGE_NAME}"
    fi
done